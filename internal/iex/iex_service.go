package iex

import (
	"context"
	"database/sql"
	"div-dash/internal/db"
	"div-dash/internal/job"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type IEXService struct {
	client     *resty.Client
	queries    *db.Queries
	db         *sql.DB
	jobService job.IJobService
}

func New(queries *db.Queries, db *sql.DB, jobService job.IJobService) *IEXService {
	client := resty.New()
	return &IEXService{client, queries, db, jobService}
}

func (i *IEXService) GetPrice(asset db.Asset) (float64, error) {

	exchanges, err := i.queries.GetExchangesOfAsset(context.Background(), asset.AssetName)
	if err != nil {
		return -1.0, err
	}

	exchange := exchanges[0]

	token := "pk_f63a9516a1d14334bcf987d1dd52af64"

	resp, err := i.client.R().
		SetQueryParam("token", token).
		SetPathParam("symbol", asset.AssetName+exchange.ExchangeSuffix).
		Get("https://cloud.iexapis.com/stable/stock/{symbol}/quote/latestPrice")

	if err != nil {
		return -1, err
	}

	body := string(resp.Body())
	if resp.StatusCode() != http.StatusOK {
		errorMsg := fmt.Sprintf("iex/GetPrice: could not get price for '%s': %s", asset.AssetName, body)
		return -1, errors.New(errorMsg)
	}

	price, err := strconv.ParseFloat(body, 64)
	if err != nil {
		return -1, err
	}

	return price, nil
}
