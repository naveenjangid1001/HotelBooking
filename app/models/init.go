package models

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v2"
)

//global database references
var db *sql.DB
var Dbmap *gorp.DbMap

//database settings
var dbName = "OnlineHotel"
var dbUser = "root"
var dbPass = "password"

//creating database connection
func Init_DB() {
	var err error

	db, err = sql.Open("mysql", dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName+"?parseTime=true")
	Dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	if err != nil {
		log.Println("Failed to connect ot database")
		log.Panic(err)
	} else {
		err = db.Ping()
		if err != nil {
			log.Println("Failed to ping Database.")
			log.Panic(err)
		} else {
			log.Println("Database connected.")
		}
	}
}
