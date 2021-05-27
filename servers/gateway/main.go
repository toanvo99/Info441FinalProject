package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

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

	sqlPass := os.Getenv("MYSQL_ROOT_PASSWORD")
	if len(sqlPass) == 0 {
		log.Fatal("No MYSQL_ROOT_PASSWORD found")
	}
	// This DSN is assuming the name of our docker image is "database"
	// and that the name of our database is mysqldb.
	// THIS IS SUBJECT TO CHANGE!!!!!
	DSN := fmt.Sprintf("root:%s@tcp(database:3306)/mysqldb", sqlPass)

	db, err := sql.Open("mysql", DSN)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	mux := http.NewServeMux()
	//mux.HandleFunc("/v1/summary", handlers.SummaryHandler)

	log.Printf("server is listening at http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
