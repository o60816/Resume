package models

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitializeDB() {
	var err error
	db, err = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/resume2?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}
