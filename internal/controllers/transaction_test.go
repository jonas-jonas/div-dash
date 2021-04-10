package controllers

import (
	"div-dash/internal/db"
	"div-dash/util/testutil"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetTransaction(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"transaction_id", "symbol", "type", "transaction_provider", "buy_in", "buy_in_date", "amount", "portfolio_id", "side"}).
		AddRow(1, "BTC", db.TransactionTypeCrypto, db.TransactionProviderBinance, 3497223, time.Date(2021, 4, 9, 10, 0, 0, 0, time.Now().Location()), 0.00034, 1, db.TransactionSideBuy)

	mock.ExpectQuery("^-- name: GetTransaction :one .*$").WithArgs(1).WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/portfolio/1/transaction/1")

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{
		"transactionId":1,
		"symbol":"BTC",
		"type":"crypto",
		"transactionProvider":"binance",
		"buyIn":34972.23,
		"buyInDate":
		"2021-04-09T10:00:00+02:00",
		"amount":"0.00034",
		"portfolioId":1,
		"side": "buy"
	}`, w.Body.String())
}

func TestGetTransactionStringId(t *testing.T) {
	_, cleanup, router := NewApi()

	defer cleanup()

	w := PerformAuthenticatedRequest(router, "GET", "/api/portfolio/1/transaction/astring")

	assert.Equal(t, 400, w.Code)
	AssertErrorObject(t, "Invalid transaction id", 400, w.Body)
}

func TestGetTransactionDbError(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: GetTransaction :one .*$").WithArgs(1).WillReturnError(errors.New("test-error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/portfolio/1/transaction/1")

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}

func TestPostTransaction(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"transaction_id", "symbol", "type", "transaction_provider", "buy_in", "buy_in_date", "amount", "portfolio_id", "side"}).
		AddRow(1, "BTC", db.TransactionTypeCrypto, db.TransactionProviderBinance, 3497223, time.Date(2021, 4, 9, 10, 0, 0, 0, time.Now().Location()), 0.00034, 1, db.TransactionSideBuy)

	mock.ExpectQuery("^-- name: CreateTransaction :one .*$").WithArgs("BTC", "crypto", "binance", 3497223, testutil.AnyTime{}, "0.00032", 1, "buy").WillReturnRows(rows)

	w := PerformAuthenticatedRequestWithBody(router, "POST", "/api/portfolio/1/transaction", `{
		"symbol": "BTC",
		"type": "crypto",
		"transactionProvider": "binance",
		"buyIn": 34972.23,
		"buyInDate": "2021-04-09T18:24:12+00:00",
		"amount": 0.00032,
		"side": "buy"
	}`)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{
		"amount":"0.00034",
		"buyIn":34972.23,
		"buyInDate":"2021-04-09T10:00:00+02:00",
		"portfolioId":1,
		"transactionId":1,
		"transactionProvider":"binance",
		"symbol":"BTC",
		"type":"crypto",
		"side": "buy"
	}`, w.Body.String())
}

func TestPostTransactionMissingField(t *testing.T) {
	_, cleanup, router := NewApi()

	defer cleanup()

	w := PerformAuthenticatedRequestWithBody(router, "POST", "/api/portfolio/1/transaction", `{
		"symbol": "BTC",
		"type": "crypto",
		"transactionProvider": "binance",
		"buyIn": 34972.23,
		"buyInDate": "2021-04-09T18:24:12+00:00"
	}`)

	assert.Equal(t, 400, w.Code)
	AssertErrorObject(t, "Key: 'createTransactionRequest.Amount' Error:Field validation for 'Amount' failed on the 'required' tag\nKey: 'createTransactionRequest.Side' Error:Field validation for 'Side' failed on the 'required' tag", 400, w.Body)
}

func TestPostTransactionStringPortfolioId(t *testing.T) {
	_, cleanup, router := NewApi()

	defer cleanup()

	w := PerformAuthenticatedRequestWithBody(router, "POST", "/api/portfolio/a-string/transaction", `{
		"symbol": "BTC",
		"type": "crypto",
		"transactionProvider": "binance",
		"buyIn": 34972.23,
		"buyInDate": "2021-04-09T18:24:12+00:00",
		"amount": 0.00032
	}`)

	assert.Equal(t, 400, w.Code)
	AssertErrorObject(t, "Invalid portfolio id", 400, w.Body)
}

func TestPostTransactionDbError(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	mock.ExpectQuery("^-- name: CreateTransaction :one .*$").WillReturnError(errors.New("test-error"))

	w := PerformAuthenticatedRequestWithBody(router, "POST", "/api/portfolio/1/transaction", `{
		"symbol": "BTC",
		"type": "crypto",
		"transactionProvider": "binance",
		"buyIn": 34972.23,
		"buyInDate": "2021-04-09T18:24:12+00:00",
		"amount": 0.00032,
		"side": "buy"
	}`)

	assert.Equal(t, 500, w.Code)
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}

func TestGetTransactions(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"transaction_id", "symbol", "type", "transaction_provider", "buy_in", "buy_in_date", "amount", "portfolio_id", "side"}).
		AddRow(1, "BTC", db.TransactionTypeCrypto, db.TransactionProviderBinance, 3497223, time.Date(2021, 4, 9, 10, 0, 0, 0, time.Now().Location()), 0.00034, 1, db.TransactionSideBuy).
		AddRow(2, "ETH", db.TransactionTypeCrypto, db.TransactionProviderBinance, 3497223, time.Date(2021, 4, 9, 10, 0, 0, 0, time.Now().Location()), 0.00034, 1, db.TransactionSideBuy).
		AddRow(3, "DOT", db.TransactionTypeCrypto, db.TransactionProviderBinance, 3497223, time.Date(2021, 4, 9, 10, 0, 0, 0, time.Now().Location()), 0.00034, 1, db.TransactionSideBuy)
	mock.ExpectQuery("^-- name: ListTransactions :many .*$").WithArgs(1).WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/portfolio/1/transaction")

	assert.Equal(t, 200, w.Code)

	var r []transactionResponse

	err := json.Unmarshal(w.Body.Bytes(), &r)

	assert.Nil(t, err)

	assert.Equal(t, 3, len(r))
}

func TestGetTransactionsStringPortfolioId(t *testing.T) {
	_, cleanup, router := NewApi()

	defer cleanup()

	w := PerformAuthenticatedRequest(router, "GET", "/api/portfolio/a-string/transaction")

	assert.Equal(t, 400, w.Code)

	AssertErrorObject(t, "Invalid portfolio id", 400, w.Body)
}

func TestGetTransactionsDbError(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	mock.ExpectQuery("^-- name: ListTransactions :many .*$").WithArgs(1).WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/portfolio/1/transaction")

	assert.Equal(t, 500, w.Code)

	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)
}
