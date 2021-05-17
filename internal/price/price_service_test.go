package price

import (
	"context"
	"div-dash/internal/db"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPriceService struct {
}

func (m *mockPriceService) GetPrice(ctx context.Context, asset db.Asset) (float64, error) {
	return 10.0, nil
}

type mockPriceServiceWithErr struct {
}

func (m *mockPriceServiceWithErr) GetPrice(ctx context.Context, asset db.Asset) (float64, error) {
	return -1, errors.New("test-price-service-error")
}

func TestGetPriceOfAsset(t *testing.T) {

	priceServices := map[string]IPriceService{
		"test-source": &mockPriceService{},
	}

	priceService := PriceService{
		priceServices: priceServices,
	}

	asset := db.Asset{
		AssetName: "test-asset",
		Type:      "crypto",
		Source:    "test-source",
	}

	ctx := context.Background()
	price, err := priceService.GetPriceOfAsset(ctx, asset)
	assert.Equal(t, price, 10.0)
	assert.Nil(t, err)
}

func TestGetPricePriceServiceErrorReturnsMinus1AndErr(t *testing.T) {

	priceServices := map[string]IPriceService{
		"test-source": &mockPriceServiceWithErr{},
	}

	priceService := PriceService{
		priceServices: priceServices,
	}
	asset := db.Asset{
		AssetName: "test-asset",
		Type:      "crypto",
		Source:    "test-source",
	}

	ctx := context.Background()
	price, err := priceService.GetPriceOfAsset(ctx, asset)
	assert.Equal(t, price, -1.0)
	assert.Equal(t, "test-price-service-error", err.Error())
}
