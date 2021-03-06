package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"Info441FinalProject/servers/gateway/handlers"
	"Info441FinalProject/servers/gateway/models"
	"Info441FinalProject/servers/gateway/sessions"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-redis/redis"
)

// Director is the director used for routing to microservices
type Director func(r *http.Request)

// CustomDirector forwards to the microservice and passes it the current user.
func CustomDirector(targets []*url.URL, ctx *handlers.HandlerContext) Director {
	var counter int32
	counter = 0
	mutex := sync.Mutex{}
	return func(r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()
		targ := targets[counter%int32(len(targets))]
		atomic.AddInt32(&counter, 1)
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Header.Del("X-User")
		sessionState := &handlers.SessionState{}
		_, err := sessions.GetState(r, ctx.SignKey, ctx.SessionStore, sessionState)
		// If there is an error, forward it to the API to deal with it.
		if err != nil {
			r.Header.Add("X-User", "{}")
		} else {
			user := sessionState.User
			userJSON, err := json.Marshal(user)
			if err != nil {
				r.Header.Add("X-User", "{}")
			} else {
				r.Header.Add("X-User", string(userJSON))
			}
		}
		r.Host = targ.Host
		r.URL.Host = targ.Host
		r.URL.Scheme = targ.Scheme
	}
}

func getURLs(addrString string) []*url.URL {
	addrsSplit := strings.Split(addrString, ",")
	URLs := make([]*url.URL, len(addrsSplit))
	for i, c := range addrsSplit {
		URL, err := url.Parse(c)
		if err != nil {
			log.Fatal(fmt.Printf("Failure to parse url %v", err))
		}
		URLs[i] = URL
	}
	return URLs
}

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")
	if len(tlsKeyPath) == 0 || len(tlsCertPath) == 0 {
		log.Printf("either TLSKEY or TLSCERT is not set")
	}
	SESSIONKEY := os.Getenv("SESSIONKEY")
	if len(SESSIONKEY) == 0 {
		log.Fatal("No SESSIONKEY found")
	}
	REDISADDR := os.Getenv("REDISADDR")
	if len(REDISADDR) == 0 {
		log.Fatal("No REDISADDR found")
	}

	sqlPass := os.Getenv("MYSQL_ROOT_PASSWORD")
	if len(sqlPass) == 0 {
		log.Fatal("No MYSQL_ROOT_PASSWORD found")
	}

	// This will be our redisStore, just keeping it unused as not sure
	// what our handler context will look like yet.
	redisStore := sessions.NewRedisStore(redis.NewClient(&redis.Options{Addr: REDISADDR,
		Password: "", DB: 0}), time.Hour)

	// This DSN is assuming the name of our docker image is "database"
	// and that the name of our database is mysqldb.
	// THIS IS SUBJECT TO CHANGE!!!!!
	DSN := fmt.Sprintf("root:%s@tcp(finaldb:3306)/finaldb", sqlPass)

	db, err := sql.Open("mysql", DSN)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	userStore := models.NewSQLStore(db)

	handlerContext := &handlers.HandlerContext{
		SignKey:      SESSIONKEY,
		TrainerStore: userStore,
		SessionStore: redisStore,
	}

	mux := http.NewServeMux()
	TEAMSSADDR := "http://teams:80"
	messagesURLs := getURLs(TEAMSSADDR)
	teamsProxy := &httputil.ReverseProxy{Director: CustomDirector(messagesURLs, handlerContext)}

	mux.HandleFunc("/v1/users", handlerContext.TrainersHandler)
	mux.HandleFunc("/v1/users/", handlerContext.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", handlerContext.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", handlerContext.SpecificSessionHandler)
	mux.Handle("/v1/", teamsProxy)
	mux.Handle("/v1/teams", teamsProxy)
	mux.Handle("/v1/teams/", teamsProxy)

	wrappedMux := handlers.NewCors(mux)
	log.Printf("server is listening at %s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))
}
