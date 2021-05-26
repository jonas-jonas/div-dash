package iex

import (
	"context"
	"div-dash/internal/db"
	"div-dash/util/mocks"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetExchanges(t *testing.T) {

	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	httpmock.RegisterResponder("GET", "https://cloud.iexapis.com/stable/ref-data/exchanges",
		httpmock.NewStringResponder(200, `[
			{
				"exchange": "test-exchange",
				"region": "de",
				"description": "A test exchange",
				"mic": "ABC",
				"exchangeSuffix": "-ABC"
			}
		]`))

	binanceService := IEXService{
		client: client,
	}

	exchanges, err := binanceService.getExchanges()

	assert.Equal(t, 1, len(exchanges))
	exchange := exchanges[0]
	assert.Equal(t, "test-exchange", exchange.Exchange)
	assert.Equal(t, "de", exchange.Region)
	assert.Equal(t, "A test exchange", exchange.Description)
	assert.Equal(t, "ABC", exchange.Mic)
	assert.Equal(t, "-ABC", exchange.ExchangeSuffix)
	assert.Nil(t, err)
}

func TestGetExchangesServerErrorReturnsNilAndError(t *testing.T) {

	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	httpmock.RegisterResponder("GET", "https://cloud.iexapis.com/stable/ref-data/exchanges",
		httpmock.NewStringResponder(500, `{"message": "error"}`))

	binanceService := IEXService{
		client: client,
	}

	exchanges, err := binanceService.getExchanges()

	assert.Nil(t, exchanges)
	assert.Equal(t, "iexService/getExchanges: error fetching exchanges: {\"message\": \"error\"}", err.Error())
}

func TestGetExchangesInvalidJSONReturnsNilAndError(t *testing.T) {

	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	httpmock.RegisterResponder("GET", "https://cloud.iexapis.com/stable/ref-data/exchanges",
		httpmock.NewStringResponder(200, `[
			{
				"exchange": "test-exchange",
				"region": "de",
				"description": "A test exchange",
				"mic": "ABC",
				"exchangeSuffix": "-ABC"`))

	binanceService := IEXService{
		client: client,
	}

	exchanges, err := binanceService.getExchanges()

	assert.Nil(t, exchanges)
	assert.Equal(t, "could not unmarshal exchanges '[\n\t\t\t{\n\t\t\t\t\"exchange\": \"test-exchange\",\n\t\t\t\t\"region\": \"de\",\n\t\t\t\t\"description\": \"A test exchange\",\n\t\t\t\t\"mic\": \"ABC\",\n\t\t\t\t\"exchangeSuffix\": \"-ABC\"': unexpected end of JSON input", err.Error())
}

func TestSaveExchanges(t *testing.T) {

	sdb, mock, _ := sqlmock.New()

	mockJobService := new(mocks.MockJobService)

	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	httpmock.RegisterResponder("GET", "https://cloud.iexapis.com/stable/ref-data/exchanges",
		httpmock.NewStringResponder(200, `[
			{
				"exchange": "test-exchange",
				"region": "de",
				"description": "A test exchange",
				"mic": "ABC",
				"exchangeSuffix": "-ABC"
			}
		]`))

	service := IEXService{
		db:         sdb,
		jobService: mockJobService,
		client:     client,
		queries:    db.New(sdb),
	}

	mockJobService.
		On("HasLastSuccessfulJobExpired", context.Background(), IEX_IMPORT_EXCHANGES_JOB_NAME, time.Duration(60*60*24*7)).
		Return(true, nil)

	mockJobService.
		On("StartJob", context.Background(), IEX_IMPORT_EXCHANGES_JOB_NAME).
		Return(db.StartJobRow{
			ID:      0,
			Started: time.Now().Unix(),
		}, nil)

	mockJobService.
		On("FinishJob", context.Background(), int32(0)).
		Return(nil)

	mock.ExpectBegin()

	mock.ExpectExec("^-- name: CreateExchange :exec.*$").
		WithArgs("test-exchange", "de", "A test exchange", "ABC", "-ABC").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	err := service.SaveExchanges()

	assert.Nil(t, err)
}

func TestSaveExchangesLastImportDidNotExpire(t *testing.T) {

	mockJobService := new(mocks.MockJobService)

	service := IEXService{
		jobService: mockJobService,
	}

	mockJobService.
		On("HasLastSuccessfulJobExpired", context.Background(), IEX_IMPORT_EXCHANGES_JOB_NAME, time.Duration(60*60*24*7)).
		Return(false, nil)

	err := service.SaveExchanges()

	assert.Equal(t, "import-iex-exchanges: last successful exchange import was less than a week ago", err.Error())
}

func TestSaveExchangesLastImportCheckFailsReturnsError(t *testing.T) {

	mockJobService := new(mocks.MockJobService)

	service := IEXService{
		jobService: mockJobService,
	}

	mockJobService.
		On("HasLastSuccessfulJobExpired", context.Background(), IEX_IMPORT_EXCHANGES_JOB_NAME, time.Duration(60*60*24*7)).
		Return(true, errors.New("test-error"))

	err := service.SaveExchanges()

	assert.Equal(t, "import-iex-exchanges: error while checking last exchange import: test-error", err.Error())
}

func TestSaveExchangesJobStartingFails(t *testing.T) {

	mockJobService := new(mocks.MockJobService)

	service := IEXService{
		jobService: mockJobService,
	}

	mockJobService.
		On("HasLastSuccessfulJobExpired", context.Background(), IEX_IMPORT_EXCHANGES_JOB_NAME, time.Duration(60*60*24*7)).
		Return(true, nil)

	mockJobService.
		On("StartJob", context.Background(), IEX_IMPORT_EXCHANGES_JOB_NAME).
		Return(db.StartJobRow{}, errors.New("test-error"))

	err := service.SaveExchanges()

	assert.Equal(t, "test-error", err.Error())
}
