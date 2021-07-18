package services

import (
	"div-dash/internal/binance"
	"div-dash/internal/config"
	"div-dash/internal/id"
	"div-dash/internal/iex"
	"div-dash/internal/job"
	"div-dash/internal/mail"
	"div-dash/internal/price"
	"div-dash/internal/token"
	"sync"

	"github.com/spf13/viper"
)

var services struct {
	TokenService   *token.TokenService
	MailService    *mail.MailService
	IdService      *id.IdService
	BinanceService *binance.BinanceService
	PriceService   *price.PriceService
	JobService     *job.JobService
	IEXService     *iex.IEXService
}

var (
	onceTokenService   sync.Once
	onceMailService    sync.Once
	onceIdService      sync.Once
	onceBinanceService sync.Once
	oncePriceService   sync.Once
	onceJobService     sync.Once
	onceIEXService     sync.Once
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

func initBinanceService() {
	services.BinanceService = binance.New(JobService(), config.DB(), config.Queries())
}

func BinanceService() *binance.BinanceService {
	onceBinanceService.Do(initBinanceService)
	return services.BinanceService
}

func initPriceService() {
	services.PriceService = price.New(BinanceService(), IEXService())
}

func PriceService() *price.PriceService {
	oncePriceService.Do(initPriceService)
	return services.PriceService
}

func initJobService() {
	services.JobService = job.New(config.Queries(), config.Logger())
}

func JobService() *job.JobService {
	onceJobService.Do(initJobService)
	return services.JobService
}

func initIEXService() {
	services.IEXService = iex.New(config.Queries(), config.DB(), JobService())
}

func IEXService() *iex.IEXService {
	onceIEXService.Do(initIEXService)
	return services.IEXService
}
