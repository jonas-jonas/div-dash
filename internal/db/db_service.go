package db

import (
	"database/sql"
	"div-dash/internal/config"
	"fmt"
	"log"
)

type (
	DBProvider interface {
		DB() *sql.DB
	}
	QueriesProvider interface {
		Queries() *Queries
	}

	dbDependencies interface {
		config.ConfigProvider
	}
	queriesDependencies interface {
		DBProvider
	}
)

func NewDB(d dbDependencies) *sql.DB {
	host := d.Config().Database.Host
	port := d.Config().Database.Port
	database := d.Config().Database.Database
	username := d.Config().Database.Username
	password := d.Config().Database.Password

	connectionString := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, database, username, password)

	log.Printf("connecting with string: %s", connectionString)
	sdb, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Print("err" + err.Error())
	}
	return sdb
}

func NewQueries(q queriesDependencies) *Queries {
	return New(q.DB())
}

func (q *Queries) GetDB() DBTX {
	return q.db
}
