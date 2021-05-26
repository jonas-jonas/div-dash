package binance

import (
	"database/sql"
	"div-dash/internal/db"
	"div-dash/internal/job"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type BinanceService struct {
	jobService *job.JobService
	db         *sql.DB
	queries    *db.Queries
	client     *resty.Client
}

func New(jobService *job.JobService, db *sql.DB, queries *db.Queries) *BinanceService {
	client := resty.New()
	return &BinanceService{
		jobService,
		db,
		queries,
		client,
	}
}

func (b *BinanceService) GetPrice(asset db.Symbol) (float64, error) {

	resp, err := b.client.R().
		SetQueryParam("symbol", asset.SymbolID+"EUR").
		Get("https://api.binance.com/api/v3/avgPrice")
	if err != nil {
		return -1, err
	}
	body := string(resp.Body())
	if resp.StatusCode() != http.StatusOK {
		return -1, fmt.Errorf("binance/GetPrice: could not get price for '%s': %s", asset.SymbolID, body)
	}
	price := gjson.Get(body, "price")

	return price.Float(), nil
}
