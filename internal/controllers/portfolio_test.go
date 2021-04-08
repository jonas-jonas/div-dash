package controllers

import (
	"div-dash/util/testutil"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetPortfolio(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"portfolio_id", "name", "user_id"}).
		AddRow(1, "Test Portfolio", 1)

	mock.ExpectQuery("^-- name: GetPortfolio :one .*$").WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/portfolio/1")

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"portfolio_id": 1, "name": "Test Portfolio", "user_id": 1}`, w.Body.String())
}

func TestGetPortfolioNoPortfolio(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{})
	mock.ExpectQuery("^-- name: GetPortfolio :one .*$").WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/portfolio/1")

	assert.Equal(t, 404, w.Code)
	assert.JSONEq(t, `{"message": "Portfolio with id '1' not found"}`, w.Body.String())
}

func TestGetPortfolioStringAsId(t *testing.T) {
	_, cleanup, router := NewApi()

	defer cleanup()

	w := PerformAuthenticatedRequest(router, "GET", "/api/portfolio/astring")

	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"message":"strconv.ParseInt: parsing \"astring\": invalid syntax", "status":500}`, w.Body.String())
}

func TestGetPortfolioDbError(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()
	mock.ExpectQuery("^-- name: GetPortfolio :one .*$").WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/portfolio/1")

	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"message": "test error", "status": 500}`, w.Body.String())
}

func TestPostPortfolio(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	rows := sqlmock.NewRows([]string{"portfolio_id", "name", "user_id"}).
		AddRow(1, "Test Portfolio", 1)

	mock.ExpectQuery("^-- name: CreatePortfolio :one .*$").WithArgs("Test Portfolio", 1).WillReturnRows(rows)

	w := PerformAuthenticatedRequestWithBody(router, "POST", "/api/portfolio", `{"name": "Test Portfolio"}`)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"portfolio_id": 1, "name": "Test Portfolio", "user_id": 1}`, w.Body.String())
}

func TestPostPortfolioDbError(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: CreatePortfolio :one .*$").WithArgs("Test Portfolio", 1).WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequestWithBody(router, "POST", "/api/portfolio", `{"name": "Test Portfolio"}`)

	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"message":"test error", "status":500}`, w.Body.String())
}

func TestPostPortfolioMissingField(t *testing.T) {
	_, cleanup, router := NewApi()

	defer cleanup()

	w := PerformAuthenticatedRequestWithBody(router, "POST", "/api/portfolio", `{}`)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"Key: 'createPortfolioRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`, w.Body.String())
}

func TestPutPortfolio(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	rows := sqlmock.NewRows([]string{"portfolio_id", "name", "user_id"}).
		AddRow(1, "New Test Portfolio", 1)

	mock.ExpectQuery("^-- name: UpdatePortfolio :one .*$").WithArgs(1, "New Test Portfolio").WillReturnRows(rows)

	w := PerformAuthenticatedRequestWithBody(router, "PUT", "/api/portfolio/1", `{"name": "New Test Portfolio"}`)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"portfolio_id": 1, "name": "New Test Portfolio", "user_id": 1}`, w.Body.String())
}

func TestPutPortfolioDbError(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: UpdatePortfolio :one .*$").WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequestWithBody(router, "PUT", "/api/portfolio/1", `{"name": "New Test Portfolio"}`)

	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"message":"test error", "status":500}`, w.Body.String())
}

func TestPutPortfolioStringId(t *testing.T) {

	_, cleanup, router := NewApi()

	defer cleanup()

	w := PerformAuthenticatedRequestWithBody(router, "PUT", "/api/portfolio/asd", `{"name": "New Test Portfolio"}`)

	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"message":"strconv.ParseInt: parsing \"asd\": invalid syntax", "status":500}`, w.Body.String())
}

func TestPutPortfolioMissingField(t *testing.T) {

	_, cleanup, router := NewApi()

	defer cleanup()

	w := PerformAuthenticatedRequestWithBody(router, "PUT", "/api/portfolio/1", `{}`)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"Key: 'updatePortfolioRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`, w.Body.String())
}

func TestDeletePortfolio(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectExec("^-- name: DeletePortfolio :exec .*$").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	w := PerformAuthenticatedRequest(router, "DELETE", "/api/portfolio/1")

	assert.Equal(t, 200, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestDeletePortfolioDbError(t *testing.T) {
	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectExec("^-- name: DeletePortfolio :exec .*$").WithArgs(1).WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequest(router, "DELETE", "/api/portfolio/1")

	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"message":"test error", "status":500}`, w.Body.String())
}

func TestDeletePortfolioStringId(t *testing.T) {
	_, cleanup, router := NewApi()

	defer cleanup()

	w := PerformAuthenticatedRequest(router, "DELETE", "/api/portfolio/asd")

	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"message":"strconv.ParseInt: parsing \"asd\": invalid syntax", "status":500}`, w.Body.String())
}

func TestGetPortfolios(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	rows := sqlmock.NewRows([]string{"portfolio_id", "name", "user_id"}).
		AddRow(1, "Test Portfolio 1", testutil.TestUserID).
		AddRow(2, "Test Portfolio 2", testutil.TestUserID).
		AddRow(3, "Test Portfolio 3", testutil.TestUserID)

	mock.ExpectQuery("^-- name: ListPortfolios :many .*$").WithArgs(testutil.TestUserID).WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/portfolio")

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `[{"portfolio_id":1,"name":"Test Portfolio 1","user_id":1},{"portfolio_id":2,"name":"Test Portfolio 2","user_id":1},{"portfolio_id":3,"name":"Test Portfolio 3","user_id":1}]`, w.Body.String())
}

func TestGetPortfoliosDbError(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: ListPortfolios :many .*$").WithArgs(testutil.TestUserID).WillReturnError(errors.New("test error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/portfolio")

	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"message":"test error", "status":500}`, w.Body.String())
}
