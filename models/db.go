package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}

func InitializeDB() {
	var err error
	var (
		connectionName = mustGetenv("CLOUDSQL_CONNECTION_NAME")
		user           = mustGetenv("CLOUDSQL_USER")
		password       = os.Getenv("CLOUDSQL_PASSWORD") // NOTE: password may be empty
	)

	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@cloudsql(%s)/resume?parseTime=true", user, password, connectionName))
	if err != nil {
		log.Fatalf("Could not open db: %v", err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}
