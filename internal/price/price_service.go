package price

import (
	"div-dash/internal/binance"
	"div-dash/internal/db"
	"div-dash/internal/iex"
	"time"

	"zgo.at/zcache"
)

type PriceService struct {
	cache         *zcache.Cache
	priceServices map[string]IPriceService
}

type IPriceService interface {
	GetPrice(asset db.Asset) (float64, error)
}

func New(binance *binance.BinanceService, iex *iex.IEXService) *PriceService {
	priceServices := map[string]IPriceService{
		"binance": binance,
		"iex":     iex,
	}
	cache := zcache.New(5*time.Minute, 10*time.Minute)
	return &PriceService{cache, priceServices}
}

func (p *PriceService) GetPriceOfAsset(asset db.Asset) (float64, error) {

	cacheKey := asset.Source + "/" + asset.AssetName

	if price, found := p.cache.Get(cacheKey); found {
		return price.(float64), nil
	}

	priceService := p.priceServices[asset.Source]

	price, err := priceService.GetPrice(asset)
	if err != nil {
		return -1, err
	}
	p.cache.Set(cacheKey, price, zcache.DefaultExpiration)

	return price, nil
}
