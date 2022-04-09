package services

import (
	"div-dash/internal/coingecko"
	"div-dash/internal/config"
	"div-dash/internal/id"
	"div-dash/internal/iex"
	"div-dash/internal/mail"
	"div-dash/internal/price"
	"div-dash/internal/token"
	"sync"

	"github.com/spf13/viper"
)

var services struct {
	TokenService     *token.TokenService
	MailService      *mail.MailService
	IdService        *id.IdService
	PriceService     *price.PriceService
	IEXService       *iex.IEXService
	CoingeckoService *coingecko.CoingeckoService
}

var (
	onceTokenService     sync.Once
	onceMailService      sync.Once
	onceIdService        sync.Once
	oncePriceService     sync.Once
	onceIEXService       sync.Once
	onceCoingeckoService sync.Once
)

func initTokenService() {
	services.TokenService = token.NewTokenService(config.Config().Token)
}

func TokenService() *token.TokenService {
	onceTokenService.Do(initTokenService)

	return services.TokenService
}

func initMailService() {
	password := viper.GetString("smtp.password")
	port := viper.GetInt("smtp.port")
	server := viper.GetString("smtp.server")
	services.MailService = mail.NewMailService(password, server, port)
}

func MailService() *mail.MailService {
	onceMailService.Do(initMailService)
	return services.MailService
}

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
