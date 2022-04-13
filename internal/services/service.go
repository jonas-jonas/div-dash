package services

import (
	"div-dash/internal/coingecko"
	"div-dash/internal/config"
	"div-dash/internal/id"
	"div-dash/internal/iex"
	"div-dash/internal/price"
	"sync"

	"go.uber.org/zap"
)

var services struct {
	IdService        *id.IdService
	PriceService     *price.PriceService
	IEXService       *iex.IEXService
	CoingeckoService *coingecko.CoingeckoService

	logger *zap.Logger
}

var (
	onceIdService        sync.Once
	oncePriceService     sync.Once
	onceIEXService       sync.Once
	onceCoingeckoService sync.Once
)

func initIdService() {
	services.IdService = id.New()
}

func IdService() *id.IdService {
	onceIdService.Do(initIdService)
	return services.IdService
}

func initPriceService() {
	services.PriceService = price.New(IEXService(), CoingeckoService())
}

func PriceService() *price.PriceService {
	oncePriceService.Do(initPriceService)
	return services.PriceService
}

func initIEXService() {
	services.IEXService = iex.New(config.Queries(), config.DB(), config.Config().IEX)
}

func IEXService() *iex.IEXService {
	onceIEXService.Do(initIEXService)
	return services.IEXService
}

func initCoingeckoService() {
	services.CoingeckoService = coingecko.New(config.Queries(), config.DB())
}

func CoingeckoService() *coingecko.CoingeckoService {
	onceCoingeckoService.Do(initCoingeckoService)
	return services.CoingeckoService
}
