package config

import (
	"database/sql"
	"div-dash/internal/db"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB
var queries *db.Queries = new(db.Queries)

func InitDB() {
	file := Config().Database.File
	cache := Config().Database.Cache
	mode := Config().Database.Mode
	loc := Config().Database.Loc

	source := fmt.Sprintf("file:%s?cache=%s&mode=%s&_loc=%s&parseTime=true", file, cache, mode, loc)
	Logger().Printf("connecting with string: %s", source)
	sdb, err := sql.Open("sqlite3", source)
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
