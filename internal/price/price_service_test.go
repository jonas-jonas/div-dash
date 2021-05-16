package price

import (
	"context"
	"div-dash/internal/db"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

func TestGetPrice(t *testing.T) {

	sdb, mock, _ := sqlmock.New()

	priceServices := map[string]IPriceService{
		"test-source": &mockPriceService{},
	}

	priceService := PriceService{
		priceServices: priceServices,
		queries:       db.New(sdb),
	}

	rows := sqlmock.NewRows([]string{"asset_name", "type", "source", "precision"}).
		AddRow("test-asset", "crypto", "test-source", 8)

	mock.ExpectQuery("^-- name: GetAsset :one .*$").WithArgs("test-asset").
		WillReturnRows(rows)

	ctx := context.Background()
	price, err := priceService.GetPrice(ctx, "test-asset")
	assert.Equal(t, price, 10.0)
	assert.Nil(t, err)
}

func TestGetPriceDbErrorReturnsMinus1AndErr(t *testing.T) {

	sdb, mock, _ := sqlmock.New()

	priceServices := map[string]IPriceService{
		"test-source": &mockPriceService{},
	}

	priceService := PriceService{
		priceServices: priceServices,
		queries:       db.New(sdb),
	}

	mock.ExpectQuery("^-- name: GetAsset :one .*$").WithArgs("test-asset").
		WillReturnError(errors.New("test-error"))

	ctx := context.Background()
	price, err := priceService.GetPrice(ctx, "test-asset")
	assert.Equal(t, price, -1.0)
	assert.Equal(t, "test-error", err.Error())
}

func TestGetPricePriceServiceErrorReturnsMinus1AndErr(t *testing.T) {

	sdb, mock, _ := sqlmock.New()

	priceServices := map[string]IPriceService{
		"test-source": &mockPriceServiceWithErr{},
	}

	priceService := PriceService{
		priceServices: priceServices,
		queries:       db.New(sdb),
	}

	rows := sqlmock.NewRows([]string{"asset_name", "type", "source", "precision"}).
		AddRow("test-asset", "crypto", "test-source", 8)

	mock.ExpectQuery("^-- name: GetAsset :one .*$").WithArgs("test-asset").
		WillReturnRows(rows)

	ctx := context.Background()
	price, err := priceService.GetPrice(ctx, "test-asset")
	assert.Equal(t, price, -1.0)
	assert.Equal(t, "test-price-service-error", err.Error())
}
