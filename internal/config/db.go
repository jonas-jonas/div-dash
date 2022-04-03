package config

import (
	"database/sql"
	"div-dash/internal/db"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var database *sql.DB
var queries *db.Queries = new(db.Queries)

func InitDB() {
	dbConfig := Config().Database
	host := dbConfig.Host
	port := dbConfig.Port
	database := dbConfig.Database
	username := dbConfig.Username
	password := dbConfig.Password

	connectionString := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, database, username, password)

	log.Printf("connecting with string: %s", connectionString)
	sdb, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Print("err" + err.Error())
	}
	SetDB(sdb)
}

func SetDB(sdb *sql.DB) {
	database = sdb
	queries = db.New(sdb)
}

func DB() *sql.DB {
	return database
}

func Queries() *db.Queries {
	return queries
}
