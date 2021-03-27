package config

import (
	"database/sql"
	"div-dash/internal/db"
	"fmt"

	_ "github.com/lib/pq"
)

var database *sql.DB
var queries *db.Queries = new(db.Queries)

func InitDB() {
	//TODO: Pass a config object here?
	pgUser := "postgres"
	pgPassword := "postgres"
	pgHost := "127.0.0.1"
	pgPort := 5432
	pgDB := "postgres"
	source := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgPort, pgDB)
	sdb, err := sql.Open("postgres", source)
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
