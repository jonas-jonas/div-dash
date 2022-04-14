package coingecko

import (
	"div-dash/internal/db"

	"github.com/go-resty/resty/v2"
)

type (
	coingeckoServiceDependencies interface {
		db.QueriesProvider
		db.DBProvider
	}
	CoingeckoServiceProvider interface {
		CoingeckoService() *CoingeckoService
	}

	CoingeckoService struct {
		coingeckoServiceDependencies
		client *resty.Client
	}
)

func New(c coingeckoServiceDependencies) *CoingeckoService {
	client := resty.New()
	client.SetHostURL("https://api.coingecko.com/api/v3")
	return &CoingeckoService{
		coingeckoServiceDependencies: c,
		client:                       client,
	}
}
