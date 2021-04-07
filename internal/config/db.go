package config

import (
	"database/sql"
	"div-dash/internal/db"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var database *sql.DB
var queries *db.Queries = new(db.Queries)

func InitDB() {
	pgUser := viper.GetString("database.username")
	pgPassword := viper.GetString("database.password")
	pgHost := viper.GetString("database.host")
	pgPort := viper.GetInt("database.port")
	pgDB := viper.GetString("database.database")
	pgSslMode := viper.GetString("database.sslmode")
	source := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", pgUser, pgPassword, pgHost, pgPort, pgDB, pgSslMode)
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
