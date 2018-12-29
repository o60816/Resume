package models

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitializeDB() {
	var err error
	db, err = sql.Open("mysql", "root:12345678@tcp(35.234.38.117:3306)/resume?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}
