package controllers

import (
	"div-dash/util/testutil"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
	mock, cleanup, router := NewApi()

	// TODO: Mock Price Services here

	defer cleanup()
	rows := sqlmock.NewRows([]string{"symbol", "cost_basis", "amount"}).
		AddRow("BTC", float64(10000), float64(20))

	mock.ExpectQuery("^-- name: GetBalance :many .*$").
		WithArgs(testutil.TestUserID).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"asset_name", "type", "source", "precision"}).
		AddRow("BTC", "crypto", "binance", 8)

	mock.ExpectQuery("^-- name: GetAsset :one .*$").
		WithArgs("BTC").
		WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/balance")

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `["amount":20, "asset":map[string]interface {}{"assetName":"BTC", "precision":8, "source":"binance", "type":"crypto"}, "costBasis":100, "fiatAssetPrice":35900.4684846, "fiatValue":718009.3696920001, "plAbsolute":717909.3696920001, "plPercent":7179.093696920001}]`, w.Body.String())
}

func TestGetBalanceDbErrorOnBalance(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: GetBalance :many .*$").
		WithArgs(testutil.TestUserID).
		WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/balance")

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)

}
