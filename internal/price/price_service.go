package price

import (
	"context"
	"div-dash/internal/db"
	"time"

	"zgo.at/zcache"
)

type (
	priceServiceDependencies interface{}

	PriceServiceProvider interface {
		PriceService() *PriceService
	}

	PriceService interface {
		GetPrice(ctx context.Context, asset db.Symbol) (float64, error)
	}

	priceService struct {
		priceServiceDependencies
		cache *zcache.Cache
	}
)

func New(p priceServiceDependencies) PriceService {
	cache := zcache.New(5*time.Minute, 10*time.Minute)
	return &priceService{
		priceServiceDependencies: p,
		cache:                    cache,
	}
}

func (p *priceService) GetPrice(ctx context.Context, symbol db.Symbol) (float64, error) {

	return 0.0, nil
}
