package controllers

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestSearchSymbol(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()
	rows := sqlmock.NewRows([]string{"symbol_id", "type", "source", "precision", "symbol_name"}).
		AddRow("test", "test-type", "test-source", 8, "test-name").
		AddRow("test1", "test-type", "test-source2", 8, "test-name1").
		AddRow("test2", "test-type", "test-source4", 8, "test-name2")

	mock.ExpectQuery("^-- name: SearchSymbol :many .*$").WithArgs("%test%", 3).WillReturnRows(rows)

	w := PerformAuthenticatedRequest(router, "GET", "/api/symbol/search?query=test&count=3")

	assert.Equal(t, 200, w.Code)

	assert.Nil(t, mock.ExpectationsWereMet())
	json := gjson.Parse(w.Body.String())

	assert.Equal(t, 3.0, json.Get("#").Num)

	assert.Equal(t, "test", json.Get("0.symbolID").Str)
	assert.Equal(t, "test1", json.Get("1.symbolID").Str)
	assert.Equal(t, "test2", json.Get("2.symbolID").Str)
	assert.Equal(t, "test-type", json.Get("0.type").Str)
	assert.Equal(t, "test-source", json.Get("0.source").Str)
	assert.Equal(t, 8.0, json.Get("0.precision").Num)
	assert.Equal(t, "test-name", json.Get("0.symbolName").Str)

}

func TestSearchSymbolCountIsStringReturnsError(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	w := PerformAuthenticatedRequest(router, "GET", "/api/symbol/search?query=test&count=not-a-number")

	assert.Equal(t, 400, w.Code)

	assert.Nil(t, mock.ExpectationsWereMet())
	AssertErrorObject(t, "param 'count' must be int", 400, w.Body)

}

func TestSearchSymbolDatabaseError(t *testing.T) {

	mock, cleanup, router := NewApi()

	defer cleanup()

	mock.ExpectQuery("^-- name: SearchSymbol :many .*$").WithArgs("%test%", 3).WillReturnError(errors.New("test-error"))

	w := PerformAuthenticatedRequest(router, "GET", "/api/symbol/search?query=test&count=3")

	assert.Equal(t, 500, w.Code)

	assert.Nil(t, mock.ExpectationsWereMet())
	AssertErrorObject(t, "An internal server error occured. Please try again later.", 500, w.Body)

}
