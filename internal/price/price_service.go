package price

import (
	"context"
	"div-dash/internal/coingecko"
	"div-dash/internal/db"
	"div-dash/internal/iex"
	"div-dash/internal/model"
	"fmt"
	"time"

	"zgo.at/zcache"
)

type PriceService struct {
	cache          *zcache.Cache
	priceServices  map[string]IPriceService
	detailServices map[string]IDetailService
	chartService   map[string]IChartService
}

type IPriceService interface {
	GetPrice(ctx context.Context, asset db.Symbol) (float64, error)
}

type IDetailService interface {
	GetDetails(ctx context.Context, asset db.Symbol) (model.SymbolDetails, error)
}

type IChartService interface {
	GetChart(ctx context.Context, asset db.Symbol, span int) (model.Chart, error)
}

func New(iex *iex.IEXService, coingecko *coingecko.CoingeckoService) *PriceService {
	priceServices := map[string]IPriceService{
		"iex":       iex,
		"coingecko": coingecko,
	}
	detailServices := map[string]IDetailService{
		"iex":       iex,
		"coingecko": coingecko,
	}
	chartServices := map[string]IChartService{
		"iex":       iex,
		"coingecko": coingecko,
	}
	cache := zcache.New(5*time.Minute, 10*time.Minute)
	return &PriceService{cache, priceServices, detailServices, chartServices}
}

func (p *PriceService) GetPriceOfAsset(ctx context.Context, asset db.Symbol) (float64, error) {

	cacheKey := asset.Source + "/" + asset.SymbolID

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

func (p *PriceService) GetDetails(ctx context.Context, asset db.Symbol) (model.SymbolDetails, error) {
	detailService := p.detailServices[asset.Source]

	return detailService.GetDetails(ctx, asset)
}

func (p *PriceService) GetChart(ctx context.Context, asset db.Symbol, span int) (model.Chart, error) {
	if chartService, ok := p.chartService[asset.Source]; ok {
		return chartService.GetChart(ctx, asset, span)

	}
	return model.Chart{}, fmt.Errorf("no chart service registered for source %s", asset.Source)

}
