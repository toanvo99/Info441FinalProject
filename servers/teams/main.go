package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"Info441FinalProject/servers/teams/teamsrc"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	addr := ":80"
	if len(addr) == 0 {
		addr = ":80"
	}
	sqlPass := os.Getenv("MYSQL_ROOT_PASSWORD")
	if len(sqlPass) == 0 {
		log.Fatal("No MYSQL_ROOT_PASSWORD found")
	}
	DSN := fmt.Sprintf("root:%s@tcp(database:3306)/mysqldb?parseTime=true", sqlPass)
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}

	teamContext := &teamsrc.TeamContext{
		TeamStore: teamsrc.NewTeamSQLStore(db),
	}

	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/", teamContext.TeamManageHandler)
	mux.HandleFunc("/v1/teams/", teamContext.TeamBuilderHandler)
	mux.HandleFunc("/v1/teams", teamContext.AllTeamHandler)

	log.Printf("server is listening at http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
