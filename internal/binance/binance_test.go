package binance

import (
	"context"
	"div-dash/internal/db"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetPrice(t *testing.T) {

	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	httpmock.RegisterResponder("GET", "https://api.binance.com/api/v3/avgPrice?symbol=BTCEUR",
		httpmock.NewStringResponder(200, `{"mins": 5, "price": 123.987}`))

	binanceService := BinanceService{
		client: client,
	}

	ctx := context.Background()

	price, err := binanceService.GetPrice(ctx, db.Asset{
		AssetName: "BTC",
	})

	assert.Equal(t, 123.987, price)
	assert.Nil(t, err)
}

func TestGetPriceErrorResponseReturnMinus1AndError(t *testing.T) {

	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	httpmock.RegisterResponder("GET", "https://api.binance.com/api/v3/avgPrice?symbol=BTCEUR",
		httpmock.NewStringResponder(500, `{"error": "test-error"}`))

	binanceService := BinanceService{
		client: client,
	}

	ctx := context.Background()

	price, err := binanceService.GetPrice(ctx, db.Asset{
		AssetName: "BTC",
	})

	assert.Equal(t, -1.0, price)
	assert.Equal(t, "binance/GetPrice: could not get price for 'BTC': {\"error\": \"test-error\"}", err.Error())
}
