package binance

import (
	"context"
	"database/sql"
	"div-dash/internal/config"
	"div-dash/internal/db"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"
)

type ExchangeInfo struct {
	Timezone        string        `json:"timezone"`
	Servertime      int64         `json:"serverTime"`
	Ratelimits      []Ratelimits  `json:"rateLimits"`
	Exchangefilters []interface{} `json:"exchangeFilters"`
	Symbols         []Symbols     `json:"symbols"`
}
type Ratelimits struct {
	Ratelimittype string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	Intervalnum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}
type Filters struct {
	Filtertype       string `json:"filterType"`
	Minprice         string `json:"minPrice,omitempty"`
	Maxprice         string `json:"maxPrice,omitempty"`
	Ticksize         string `json:"tickSize,omitempty"`
	Multiplierup     string `json:"multiplierUp,omitempty"`
	Multiplierdown   string `json:"multiplierDown,omitempty"`
	Avgpricemins     int    `json:"avgPriceMins,omitempty"`
	Minqty           string `json:"minQty,omitempty"`
	Maxqty           string `json:"maxQty,omitempty"`
	Stepsize         string `json:"stepSize,omitempty"`
	Minnotional      string `json:"minNotional,omitempty"`
	Applytomarket    bool   `json:"applyToMarket,omitempty"`
	Limit            int    `json:"limit,omitempty"`
	Maxnumalgoorders int    `json:"maxNumAlgoOrders,omitempty"`
	Maxnumorders     int    `json:"maxNumOrders,omitempty"`
}
type Symbols struct {
	Symbol                     string    `json:"symbol"`
	Status                     string    `json:"status"`
	Baseasset                  string    `json:"baseAsset"`
	Baseassetprecision         int32     `json:"baseAssetPrecision"`
	Quoteasset                 string    `json:"quoteAsset"`
	Quoteprecision             int       `json:"quotePrecision"`
	Quoteassetprecision        int       `json:"quoteAssetPrecision"`
	Basecommissionprecision    int       `json:"baseCommissionPrecision"`
	Quotecommissionprecision   int       `json:"quoteCommissionPrecision"`
	Ordertypes                 []string  `json:"orderTypes"`
	Icebergallowed             bool      `json:"icebergAllowed"`
	Ocoallowed                 bool      `json:"ocoAllowed"`
	Quoteorderqtymarketallowed bool      `json:"quoteOrderQtyMarketAllowed"`
	Isspottradingallowed       bool      `json:"isSpotTradingAllowed"`
	Ismargintradingallowed     bool      `json:"isMarginTradingAllowed"`
	Filters                    []Filters `json:"filters"`
	Permissions                []string  `json:"permissions"`
}

const IMPORT_JOB_NAME = "import-binance-exchange-info"

func (b *BinanceService) getExchangeInfo() (*ExchangeInfo, error) {

	file, err := ioutil.ReadFile("data/binance/exchangeInfo.json")
	if err != nil {
		return nil, err
	}

	var exchangeInfo ExchangeInfo

	err = json.Unmarshal([]byte(file), &exchangeInfo)
	if err != nil {
		return nil, err
	}
	return &exchangeInfo, nil
}

func (b *BinanceService) SaveExchangeInfo() error {
	ctx := context.Background()
	week := 60 * 60 * 24 * 7
	expired, err := b.jobService.HasLastSuccessfulJobExpired(ctx, IMPORT_JOB_NAME, time.Duration(week))
	if !expired || err == nil {
		return errors.New("binance-exchange-info-import: last successful import was less than a week ago")
	}

	exchangeInfo, err := b.getExchangeInfo()
	if err != nil {
		return err
	}
	count := len(exchangeInfo.Symbols)
	config.Logger().Printf("Importing %v Binance Assets...", count)
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	queries := b.queries.WithTx(tx)

	job, err := b.jobService.StartJob(ctx, IMPORT_JOB_NAME)
	if err != nil {
		return err
	}
	for _, symbol := range exchangeInfo.Symbols {
		err := queries.AddAsset(ctx, db.AddAssetParams{
			AssetName: symbol.Baseasset,
			Type:      "crypto",
			Source:    "binance",
			Precision: sql.NullInt32{
				Int32: symbol.Baseassetprecision,
				Valid: true,
			},
		})
		if err != nil {
			config.Logger().Printf("Could not import symbol %s because %s, ignoring...", symbol.Baseasset, err.Error())
		}
	}
	err = tx.Commit()
	if err != nil {
		b.jobService.FailJob(ctx, job.ID, err.Error())
	}
	config.Logger().Printf("Imported %v Binance Assets successfully.", count)
	b.jobService.FinishJob(ctx, job.ID)
	return nil
}
