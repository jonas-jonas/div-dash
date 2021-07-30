package price

import (
	"div-dash/internal/db"
	"div-dash/internal/iex"
	"div-dash/internal/model"
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
	GetPrice(asset db.Symbol) (float64, error)
}

type IDetailService interface {
	GetDetails(asset db.Symbol) (model.SymbolDetails, error)
}

type IChartService interface {
	GetChart(asset db.Symbol, span int) (model.Chart, error)
}

func New(iex *iex.IEXService) *PriceService {
	priceServices := map[string]IPriceService{
		"iex": iex,
	}
	detailServices := map[string]IDetailService{
		"iex": iex,
	}
	chartServices := map[string]IChartService{
		"iex": iex,
	}
	cache := zcache.New(5*time.Minute, 10*time.Minute)
	return &PriceService{cache, priceServices, detailServices, chartServices}
}

func (p *PriceService) GetPriceOfAsset(asset db.Symbol) (float64, error) {

	cacheKey := asset.Source + "/" + asset.SymbolID

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

func (p *PriceService) GetDetails(asset db.Symbol) (model.SymbolDetails, error) {
	detailService := p.detailServices[asset.Source]

	return detailService.GetDetails(asset)
}

func (p *PriceService) GetChart(asset db.Symbol, span int) (model.Chart, error) {
	chartService := p.chartService[asset.Source]

	return chartService.GetChart(asset, span)
}
