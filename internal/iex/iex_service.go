package iex

import (
	"context"
	"database/sql"
	"div-dash/internal/db"
	"div-dash/internal/job"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"zgo.at/zcache"
)

type IEXService struct {
	client     *resty.Client
	queries    *db.Queries
	db         *sql.DB
	jobService job.IJobService
	quoteCache *zcache.Cache
}

func New(queries *db.Queries, db *sql.DB, jobService job.IJobService) *IEXService {
	client := resty.New()
	quoteCache := zcache.New(zcache.NoExpiration, -1)
	return &IEXService{client, queries, db, jobService, quoteCache}
}

var exchangeWeights = map[string]int{
	"GY": 10,
}

type Quote struct {
	Symbol                 string  `json:"symbol"`
	CompanyName            string  `json:"companyName"`
	PrimaryExchange        string  `json:"primaryExchange"`
	CalculationPrice       string  `json:"calculationPrice"`
	Open                   float64 `json:"open"`
	OpenTime               int64   `json:"openTime"`
	OpenSource             string  `json:"openSource"`
	Close                  float64 `json:"close"`
	CloseTime              int64   `json:"closeTime"`
	CloseSource            string  `json:"closeSource"`
	High                   float64 `json:"high"`
	HighTime               int64   `json:"highTime"`
	HighSource             string  `json:"highSource"`
	Low                    float64 `json:"low"`
	LowTime                int64   `json:"lowTime"`
	LowSource              string  `json:"lowSource"`
	LatestPrice            float64 `json:"latestPrice"`
	LatestSource           string  `json:"latestSource"`
	LatestTime             string  `json:"latestTime"`
	LatestUpdate           int64   `json:"latestUpdate"`
	LatestVolume           int     `json:"latestVolume"`
	IexRealtimePrice       float64 `json:"iexRealtimePrice"`
	IexRealtimeSize        int     `json:"iexRealtimeSize"`
	IexLastUpdated         int64   `json:"iexLastUpdated"`
	DelayedPrice           float64 `json:"delayedPrice"`
	DelayedPriceTime       int64   `json:"delayedPriceTime"`
	OddLotDelayedPrice     float64 `json:"oddLotDelayedPrice"`
	OddLotDelayedPriceTime int64   `json:"oddLotDelayedPriceTime"`
	ExtendedPrice          float64 `json:"extendedPrice"`
	ExtendedChange         float64 `json:"extendedChange"`
	ExtendedChangePercent  float64 `json:"extendedChangePercent"`
	ExtendedPriceTime      int64   `json:"extendedPriceTime"`
	PreviousClose          float64 `json:"previousClose"`
	PreviousVolume         int     `json:"previousVolume"`
	Change                 float64 `json:"change"`
	ChangePercent          float64 `json:"changePercent"`
	Volume                 int     `json:"volume"`
	IexMarketPercent       float64 `json:"iexMarketPercent"`
	IexVolume              int     `json:"iexVolume"`
	AvgTotalVolume         int     `json:"avgTotalVolume"`
	IexBidPrice            int     `json:"iexBidPrice"`
	IexBidSize             int     `json:"iexBidSize"`
	IexAskPrice            int     `json:"iexAskPrice"`
	IexAskSize             int     `json:"iexAskSize"`
	IexOpen                float64 `json:"iexOpen"`
	IexOpenTime            int64   `json:"iexOpenTime"`
	IexClose               float64 `json:"iexClose"`
	IexCloseTime           int64   `json:"iexCloseTime"`
	MarketCap              int64   `json:"marketCap"`
	PeRatio                float64 `json:"peRatio"`
	Week52High             float64 `json:"week52High"`
	Week52Low              float64 `json:"week52Low"`
	YtdChange              float64 `json:"ytdChange"`
	LastTradeTime          int64   `json:"lastTradeTime"`
	IsUSMarketOpen         bool    `json:"isUSMarketOpen"`
}

func (i *IEXService) GetPrice(asset db.Symbol) (float64, error) {

	exchanges, err := i.queries.GetExchangesOfAsset(context.Background(), asset.SymbolID)
	if err != nil {
		return -1.0, err
	}

	var exchange db.Exchange
	lastExchangeWeight := -1

	for _, ex := range exchanges {
		weight := exchangeWeights[ex.Exchange]
		if weight > lastExchangeWeight {
			exchange = ex
		}
	}

	if quote, found := i.quoteCache.Get(asset.SymbolID + exchange.ExchangeSuffix); found {
		return quote.(Quote).LatestPrice, nil
	}

	token := "pk_f63a9516a1d14334bcf987d1dd52af64"

	resp, err := i.client.R().
		SetQueryParam("token", token).
		SetPathParam("symbol", asset.SymbolID+exchange.ExchangeSuffix).
		Get("https://cloud.iexapis.com/stable/stock/{symbol}/quote")

	if err != nil {
		return -1, err
	}

	if resp.StatusCode() != http.StatusOK {
		errorMsg := fmt.Sprintf("iex/GetPrice: could not get price for '%s': %s", asset.SymbolID, resp.Body())
		return -1, errors.New(errorMsg)
	}

	var quote Quote
	err = json.Unmarshal(resp.Body(), &quote)
	if err != nil {
		return -1, fmt.Errorf("did not receive valid json from api: %w", err)
	}
	i.quoteCache.Set(asset.SymbolID+exchange.ExchangeSuffix, quote, time.Hour)

	return quote.LatestPrice, nil
}
