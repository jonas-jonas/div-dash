package db

import (
	"database/sql"
	"div-dash/internal/config"
	"fmt"
	"log"
)

type (
	QueriesProvider interface {
		Queries() *Queries
	}

	queriesDependencies interface {
		config.ConfigProvider
	}
)

func NewQueries(q queriesDependencies) *Queries {
	host := q.Config().Database.Host
	port := q.Config().Database.Port
	database := q.Config().Database.Database
	username := q.Config().Database.Username
	password := q.Config().Database.Password

	connectionString := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, database, username, password)

	log.Printf("connecting with string: %s", connectionString)
	sdb, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Print("err" + err.Error())
	}
	return New(sdb)
}
