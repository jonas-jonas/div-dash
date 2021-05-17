package price

import (
	"context"
	"div-dash/internal/binance"
	"div-dash/internal/db"
)

type PriceService struct {
	priceServices map[string]IPriceService
}

type IPriceService interface {
	GetPrice(ctx context.Context, asset db.Asset) (float64, error)
}

func New(binance *binance.BinanceService) *PriceService {
	priceServices := map[string]IPriceService{
		"binance": binance,
	}
	return &PriceService{priceServices}
}

func (p *PriceService) GetPriceOfAsset(ctx context.Context, asset db.Asset) (float64, error) {

	priceService := p.priceServices[asset.Source]

	price, err := priceService.GetPrice(ctx, asset)
	if err != nil {
		return -1, err
	}

	return price, nil
}
