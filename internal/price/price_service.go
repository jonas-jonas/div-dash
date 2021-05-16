package price

import (
	"context"
	"div-dash/internal/binance"
	"div-dash/internal/db"
)

type PriceService struct {
	priceServices map[string]IPriceService
	queries       *db.Queries
}

type IPriceService interface {
	GetPrice(ctx context.Context, asset db.Asset) (float64, error)
}

func New(queries *db.Queries, binance *binance.BinanceService) *PriceService {
	priceServices := map[string]IPriceService{
		"binance": binance,
	}
	return &PriceService{priceServices, queries}
}

func (p *PriceService) GetPrice(ctx context.Context, assetName string) (float64, error) {
	asset, err := p.queries.GetAsset(ctx, assetName)
	if err != nil {
		return -1, err
	}

	priceService := p.priceServices[asset.Source]

	price, err := priceService.GetPrice(ctx, asset)
	if err != nil {
		return -1, err
	}

	return price, nil
}
