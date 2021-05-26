package price

import (
	"div-dash/internal/db"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"zgo.at/zcache"
)

type mockPriceService struct {
}

func (m *mockPriceService) GetPrice(asset db.Asset) (float64, error) {
	return 10.0, nil
}

type mockPriceServiceWithErr struct {
}

func (m *mockPriceServiceWithErr) GetPrice(asset db.Asset) (float64, error) {
	return -1, errors.New("test-price-service-error")
}

func TestGetPriceOfAsset(t *testing.T) {

	priceServices := map[string]IPriceService{
		"test-source": &mockPriceService{},
	}

	priceService := PriceService{
		priceServices: priceServices,
		cache:         zcache.New(1, 10),
	}

	asset := db.Asset{
		AssetName: "test-asset",
		Type:      "crypto",
		Source:    "test-source",
	}

	price, err := priceService.GetPriceOfAsset(asset)
	assert.Equal(t, price, 10.0)
	assert.Nil(t, err)
}

func TestGetPriceOfAssetCacheHit(t *testing.T) {

	priceServices := map[string]IPriceService{
		"test-source": &mockPriceService{},
	}

	testCache := zcache.New(1*time.Minute, 10*time.Minute)

	priceService := PriceService{
		priceServices: priceServices,
		cache:         testCache,
	}

	testCache.Set("test-source/test-asset", 182473.2414, zcache.DefaultExpiration)

	asset := db.Asset{
		AssetName: "test-asset",
		Type:      "crypto",
		Source:    "test-source",
	}
	price, err := priceService.GetPriceOfAsset(asset)
	assert.Equal(t, price, 182473.2414)
	assert.Nil(t, err)
}

func TestGetPricePriceServiceErrorReturnsMinus1AndErr(t *testing.T) {

	priceServices := map[string]IPriceService{
		"test-source": &mockPriceServiceWithErr{},
	}

	priceService := PriceService{
		priceServices: priceServices,
		cache:         zcache.New(1, 10),
	}
	asset := db.Asset{
		AssetName: "test-asset",
		Type:      "crypto",
		Source:    "test-source",
	}
	price, err := priceService.GetPriceOfAsset(asset)
	assert.Equal(t, price, -1.0)
	assert.Equal(t, "test-price-service-error", err.Error())
}
