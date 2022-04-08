package iex

import (
	"context"
	"database/sql"
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"zgo.at/zcache"
)

type IEXService struct {
	client     *resty.Client
	queries    *db.Queries
	db         *sql.DB
	quoteCache *zcache.Cache
	token      string
	baseUrl    string
}

func New(queries *db.Queries, db *sql.DB, iexConfig config.IEXConfiguration) *IEXService {
	client := resty.New()
	quoteCache := zcache.New(zcache.NoExpiration, -1)
	token := iexConfig.Token
	baseUrl := iexConfig.BaseUrl
	client.SetHostURL(baseUrl)
	client.SetQueryParam("token", token)
	return &IEXService{client, queries, db, quoteCache, token, baseUrl}
}

var exchangeWeights = map[string]int{
	"ETR": 10,
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

func (i *IEXService) GetPrice(ctx context.Context, asset db.Symbol) (float64, error) {

	exchanges, err := i.queries.GetExchangesOfSymbol(context.Background(), asset.SymbolID)
	if err != nil {
		return -1.0, err
	}

	var exchange db.Exchange
	lastExchangeWeight := -1

	for _, ex := range exchanges {
		weight := exchangeWeights[ex.Exchange]
		if weight > lastExchangeWeight {
			exchange = ex
			lastExchangeWeight = weight
		}
	}

	iexSymbolId := asset.SymbolID + "-" + exchange.ExchangeSuffix

	if quote, found := i.quoteCache.Get(iexSymbolId); found {
		return quote.(Quote).LatestPrice, nil
	}

	resp, err := i.client.R().
		SetPathParam("symbol", iexSymbolId).
		Get("/stock/{symbol}/quote")

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

func (i *IEXService) getCompanyDetails(symbol string) (CompanyDetails, error) {

	resp, err := i.client.R().
		SetPathParam("symbol", symbol).
		Get("/stock/{symbol}/company")

	if err != nil {
		return CompanyDetails{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		errorMsg := fmt.Sprintf("iex/getCompanyDetails: could not get details for '%s': %s", symbol, resp.Body())
		return CompanyDetails{}, errors.New(errorMsg)
	}

	var companyDetails CompanyDetails
	err = json.Unmarshal(resp.Body(), &companyDetails)

	return companyDetails, err
}

func (i *IEXService) getCompanyKeyStats(symbol string) (CompanyKeyStats, error) {

	resp, err := i.client.R().
		SetPathParam("symbol", symbol).
		Get("/stock/{symbol}/stats")

	if err != nil {
		return CompanyKeyStats{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		errorMsg := fmt.Sprintf("iex/getCompanyKeyStats: could not get details for '%s': %s", symbol, resp.Body())
		return CompanyKeyStats{}, errors.New(errorMsg)
	}

	var companyKeyStats CompanyKeyStats
	err = json.Unmarshal(resp.Body(), &companyKeyStats)

	return companyKeyStats, err
}

func assembleTags(companyDetails CompanyDetails) []model.SymbolTag {

	tags := []model.SymbolTag{}

	tags = append(tags, model.SymbolTag{
		Label: fmt.Sprintf("%d Employees", companyDetails.Employees),
		Type:  "CHIP",
	})

	for _, tag := range companyDetails.Tags {
		tags = append(tags, model.SymbolTag{
			Label: tag,
			Type:  "CHIP",
		})
	}

	tags = append(tags, model.SymbolTag{
		Label: companyDetails.Website,
		Link:  companyDetails.Website,
		Type:  "LINK",
	})
	return tags
}

func assembleIndicators(companyKeyStats CompanyKeyStats) []model.SymbolIndicator {
	return []model.SymbolIndicator{
		{
			Label:  "Market Cap",
			Format: "$0.00 a",
			Value:  float64(companyKeyStats.Marketcap),
		},
		{
			Label:  "PE Ratio",
			Format: "0.00",
			Value:  companyKeyStats.PeRatio,
		},
		{
			Label:  "Dividend Yield",
			Format: "0.00%",
			Value:  companyKeyStats.DividendYield,
		},
		{
			Label:  "EPS",
			Format: "$0.00",
			Value:  companyKeyStats.TtmEPS,
		},
	}
}

func (i *IEXService) GetDetails(ctx context.Context, asset db.Symbol) (model.SymbolDetails, error) {

	exchanges, err := i.queries.GetExchangesOfSymbol(context.Background(), asset.SymbolID)
	if err != nil {
		return model.SymbolDetails{}, err
	}

	var exchange db.Exchange
	lastExchangeWeight := -1

	for _, ex := range exchanges {
		weight := exchangeWeights[ex.Exchange]
		if weight > lastExchangeWeight {
			exchange = ex
			lastExchangeWeight = weight
		}
	}

	iexSymbolId := asset.SymbolID + "-" + exchange.ExchangeSuffix

	companyDetails, err := i.getCompanyDetails(iexSymbolId)

	if err != nil {
		return model.SymbolDetails{}, err
	}
	companyKeyStats, err := i.getCompanyKeyStats(iexSymbolId)
	if err != nil {
		return model.SymbolDetails{}, err
	}

	return model.SymbolDetails{
		Type:        asset.Type,
		Name:        companyDetails.CompanyName,
		Description: companyDetails.Description,
		Tags:        assembleTags(companyDetails),
		Indicators:  assembleIndicators(companyKeyStats),
		Dates:       []model.SymbolDate{},
	}, nil
}

func (i *IEXService) GetChart(ctx context.Context, asset db.Symbol, span int) (model.Chart, error) {

	exchanges, err := i.queries.GetExchangesOfSymbol(context.Background(), asset.SymbolID)
	if err != nil {
		return model.Chart{}, err
	}

	var exchange db.Exchange
	lastExchangeWeight := -1

	for _, ex := range exchanges {
		weight := exchangeWeights[ex.Exchange]
		if weight > lastExchangeWeight {
			exchange = ex
			lastExchangeWeight = weight
		}
	}

	iexSymbolId := asset.SymbolID + "-" + exchange.ExchangeSuffix

	resp, err := i.client.R().
		SetPathParam("symbol", iexSymbolId).
		SetPathParam("span", strconv.Itoa(span)+"y").
		SetQueryParam("chartCloseOnly", "true").
		Get("/stock/{symbol}/chart/{span}")

	if err != nil {
		return model.Chart{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		errorMsg := fmt.Sprintf("iex/GetChart: could not get chart for '%s': %s", asset.SymbolID, resp.Body())
		return model.Chart{}, errors.New(errorMsg)
	}

	var chart []ChartEntry
	err = json.Unmarshal(resp.Body(), &chart)
	if err != nil {
		return model.Chart{}, err
	}

	chartResult := model.Chart{}

	for _, chartEntry := range chart {
		chartResult = append(chartResult, model.ChartEntry{
			Date:  chartEntry.Date,
			Price: chartEntry.Close,
		})
	}
	return chartResult, nil
}
