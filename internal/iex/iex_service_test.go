package iex

import (
	"div-dash/internal/db"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetPrice(t *testing.T) {

	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	sdb, mock, _ := sqlmock.New()
	queries := db.New(sdb)

	httpmock.RegisterResponder("GET", "https://cloud.iexapis.com/stable/stock/LUS-TEST/quote/latestPrice",
		httpmock.NewStringResponder(200, `3923.423`))

	binanceService := IEXService{
		client:  client,
		queries: queries,
	}

	rows := sqlmock.NewRows([]string{"exchange", "exchange_suffix", "region", "description", "mic"}).
		AddRow("test-exchange", "-TEST", "DE", "Test", "TEST")

	mock.ExpectQuery("^-- name: GetExchangesOfAsset :many.*$").
		WithArgs("LUS").
		WillReturnRows(rows)

	price, err := binanceService.GetPrice(db.Asset{
		AssetName: "LUS",
	})

	assert.Equal(t, 3923.423, price)
	assert.Nil(t, err)
}

func TestGetPriceForAssetWithoutPrice(t *testing.T) {

	sdb, mock, _ := sqlmock.New()
	queries := db.New(sdb)

	binanceService := IEXService{
		queries: queries,
	}

	mock.ExpectQuery("^-- name: GetExchangesOfAsset :many.*$").
		WithArgs("LUS").
		WillReturnError(errors.New("test-error"))

	price, err := binanceService.GetPrice(db.Asset{
		AssetName: "LUS",
	})

	assert.Equal(t, -1.0, price)
	assert.Equal(t, "test-error", err.Error())
}

func TestGetPriceNonOKStatusReturnsMinus1(t *testing.T) {

	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	sdb, mock, _ := sqlmock.New()
	queries := db.New(sdb)

	httpmock.RegisterResponder("GET", "https://cloud.iexapis.com/stable/stock/LUS-TEST/quote/latestPrice",
		httpmock.NewStringResponder(500, `{"msg": "error"}`))

	binanceService := IEXService{
		client:  client,
		queries: queries,
	}

	rows := sqlmock.NewRows([]string{"exchange", "exchange_suffix", "region", "description", "mic"}).
		AddRow("test-exchange", "-TEST", "DE", "Test", "TEST")

	mock.ExpectQuery("^-- name: GetExchangesOfAsset :many.*$").
		WithArgs("LUS").
		WillReturnRows(rows)

	price, err := binanceService.GetPrice(db.Asset{
		AssetName: "LUS",
	})

	assert.Equal(t, -1.0, price)
	assert.Equal(t, `iex/GetPrice: could not get price for 'LUS': {"msg": "error"}`, err.Error())
}

func TestGetPriceNonNumberReturnMinus1(t *testing.T) {

	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	sdb, mock, _ := sqlmock.New()
	queries := db.New(sdb)

	httpmock.RegisterResponder("GET", "https://cloud.iexapis.com/stable/stock/LUS-TEST/quote/latestPrice",
		httpmock.NewStringResponder(200, `not-a-number`))

	binanceService := IEXService{
		client:  client,
		queries: queries,
	}

	rows := sqlmock.NewRows([]string{"exchange", "exchange_suffix", "region", "description", "mic"}).
		AddRow("test-exchange", "-TEST", "DE", "Test", "TEST")

	mock.ExpectQuery("^-- name: GetExchangesOfAsset :many.*$").
		WithArgs("LUS").
		WillReturnRows(rows)

	price, err := binanceService.GetPrice(db.Asset{
		AssetName: "LUS",
	})

	assert.Equal(t, -1.0, price)
	assert.Equal(t, `strconv.ParseFloat: parsing "not-a-number": invalid syntax`, err.Error())
}
