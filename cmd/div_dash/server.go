package main

import (
	"div-dash/internal/coingecko"
	"div-dash/internal/config"
	"div-dash/internal/controllers"
	"div-dash/internal/iex"
	"div-dash/internal/job"
	"div-dash/internal/logging"
	"div-dash/internal/services"
	"log"
	"time"

	ginzap "github.com/gin-contrib/zap"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	// log.SetFlags(log.Lshortfile | log.LstdFlags)
	logger, err := logging.InitLogger()
	if err != nil {
		log.Panicf("Could not initialize logger")
		return
	}
	config.ReadConfig()
	config.InitDB()
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, false))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	// TODO: Replace usage of config.Queries() with a properly initiated instance
	jobService := job.NewJobService(config.Queries(), logger)

	jobService.RunJob(iex.IEXExchangesImportJob, services.IEXService().SaveExchanges)
	jobService.RunJob(iex.IEXImportSymbolsJob, services.IEXService().SaveSymbols)
	jobService.RunJob(iex.ISINAndWKNImportJob, services.IEXService().ImportISINAndWKN)
	jobService.RunJob(coingecko.CoingeckoImportCoinsJob, services.CoingeckoService().ImportCryptoSymbols)

	// TODO: Replace usage of config.Queries() with a properly initiated instance
	controllerRouter := controllers.NewControllerRouter(r, config.Queries(), logger)
	controllerRouter.RegisterRoutes()

	port := viper.GetString("server.port")
	r.Run(":" + port)
}
