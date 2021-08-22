package coingecko

import (
	"database/sql"
	"div-dash/internal/db"

	"github.com/go-resty/resty/v2"
)

type CoingeckoService struct {
	client  *resty.Client
	queries *db.Queries
	db      *sql.DB
}

func New(queries *db.Queries, db *sql.DB) *CoingeckoService {
	client := resty.New()
	client.SetHostURL("https://api.coingecko.com/api/v3")
	return &CoingeckoService{
		queries: queries,
		client:  client,
		db:      db,
	}
}
