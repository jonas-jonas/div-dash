package iex

import (
	"div-dash/internal/db"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"zgo.at/zcache"
)

const validQuoteJson = `{
	"symbol": "ADM-GY",
	"companyName": "Archer-Daniels-Midland Company",
	"primaryExchange": "XETRA",
	"calculationPrice": "close",
	"open": null,
	"openTime": null,
	"openSource": "official",
	"close": null,
	"closeTime": null,
	"closeSource": "official",
	"high": null,
	"highTime": 1622491200000,
	"highSource": "Close",
	"low": null,
	"lowTime": 1622232000000,
	"lowSource": "Close",
	"latestPrice": 54.86,
	"latestSource": "Close",
	"latestTime": "May 31, 2021",
	"latestUpdate": 1622491200000,
	"latestVolume": null,
	"iexRealtimePrice": null,
	"iexRealtimeSize": null,
	"iexLastUpdated": null,
	"delayedPrice": null,
	"delayedPriceTime": null,
	"oddLotDelayedPrice": null,
	"oddLotDelayedPriceTime": null,
	"extendedPrice": null,
	"extendedChange": null,
	"extendedChangePercent": null,
	"extendedPriceTime": null,
	"previousClose": 54.66,
	"previousVolume": 11,
	"change": 0.2,
	"changePercent": 0.00366,
	"volume": null,
	"iexMarketPercent": null,
	"iexVolume": null,
	"avgTotalVolume": 192,
	"iexBidPrice": null,
	"iexBidSize": null,
	"iexAskPrice": null,
	"iexAskSize": null,
	"iexOpen": null,
	"iexOpenTime": null,
	"iexClose": null,
	"iexCloseTime": null,
	"marketCap": 30650397755,
	"peRatio": 14.91,
	"week52High": 56.46,
	"week52Low": 33.38,
	"ytdChange": 0.3649503225806452,
	"lastTradeTime": null,
	"isUSMarketOpen": true
}`

func TestGetPrice(t *testing.T) {

	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	sdb, mock, _ := sqlmock.New()
	queries := db.New(sdb)

	httpmock.RegisterResponder("GET", "https://cloud.iexapis.com/stable/stock/LUS-TEST/quote",
		httpmock.NewStringResponder(200, validQuoteJson))

	iexService := IEXService{
		client:     client,
		queries:    queries,
		quoteCache: zcache.New(zcache.DefaultExpiration, -1),
	}

	rows := sqlmock.NewRows([]string{"exchange", "exchange_suffix", "region", "description", "mic"}).
		AddRow("test-exchange", "-TEST", "DE", "Test", "TEST")

	mock.ExpectQuery("^-- name: GetExchangesOfAsset :many.*$").
		WithArgs("LUS").
		WillReturnRows(rows)

	price, err := iexService.GetPrice(db.Symbol{
		SymbolID: "LUS",
	})

	assert.Equal(t, 54.86, price)
	assert.Nil(t, err)
}

func TestGetPriceForAssetWithoutPrice(t *testing.T) {

	sdb, mock, _ := sqlmock.New()
	queries := db.New(sdb)

	iexService := IEXService{
		queries:    queries,
		quoteCache: zcache.New(zcache.DefaultExpiration, -1),
	}

	mock.ExpectQuery("^-- name: GetExchangesOfAsset :many.*$").
		WithArgs("LUS").
		WillReturnError(errors.New("test-error"))

	price, err := iexService.GetPrice(db.Symbol{
		SymbolID: "LUS",
	})

	assert.Equal(t, -1.0, price)
	assert.Equal(t, "test-error", err.Error())
}

func TestGetPriceNonOKStatusReturnsMinus1(t *testing.T) {

	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	sdb, mock, _ := sqlmock.New()
	queries := db.New(sdb)

	httpmock.RegisterResponder("GET", "https://cloud.iexapis.com/stable/stock/LUS-TEST/quote",
		httpmock.NewStringResponder(500, `{"msg": "error"}`))

	iexService := IEXService{
		client:     client,
		queries:    queries,
		quoteCache: zcache.New(zcache.DefaultExpiration, -1),
	}

	rows := sqlmock.NewRows([]string{"exchange", "exchange_suffix", "region", "description", "mic"}).
		AddRow("test-exchange", "-TEST", "DE", "Test", "TEST")

	mock.ExpectQuery("^-- name: GetExchangesOfAsset :many.*$").
		WithArgs("LUS").
		WillReturnRows(rows)

	price, err := iexService.GetPrice(db.Symbol{
		SymbolID: "LUS",
	})

	assert.Equal(t, -1.0, price)
	assert.Equal(t, `iex/GetPrice: could not get price for 'LUS': {"msg": "error"}`, err.Error())
}
