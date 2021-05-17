package price

import (
	"context"
	"div-dash/internal/binance"
	"div-dash/internal/db"
	"time"

	"zgo.at/zcache"
)

type PriceService struct {
	cache         *zcache.Cache
	priceServices map[string]IPriceService
}

type IPriceService interface {
	GetPrice(ctx context.Context, asset db.Asset) (float64, error)
}

func New(binance *binance.BinanceService) *PriceService {
	priceServices := map[string]IPriceService{
		"binance": binance,
	}
	cache := zcache.New(5*time.Minute, 10*time.Minute)
	return &PriceService{cache, priceServices}
}

func (p *PriceService) GetPriceOfAsset(ctx context.Context, asset db.Asset) (float64, error) {

	cacheKey := asset.Source + "/" + asset.AssetName

	if price, found := p.cache.Get(cacheKey); found {
		return price.(float64), nil
	}

	priceService := p.priceServices[asset.Source]

	price, err := priceService.GetPrice(ctx, asset)
	if err != nil {
		return -1, err
	}
	p.cache.Set(cacheKey, price, zcache.DefaultExpiration)

	return price, nil
}
