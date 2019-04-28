package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	host     = "ec2-79-125-4-72.eu-west-1.compute.amazonaws.com"
	port     = 5432
	user     = "cqsuvjssroxgiw"
	password = "9f93968991a9469168ca773a4e7d170327d2311275ce548a4a7500663929ae50"
	dbname   = "d9a6lp0thaifvl"
)

func GetDB() *sql.DB {
	if db == nil {
		ConnectDB()
	}
	return db
}

func ConnectDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require\n", host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to Database!")
}

func DisconnectDB() {
	if db != nil {
		db.Close()
	}
	db = nil
}
