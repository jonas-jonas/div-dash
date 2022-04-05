package main

import (
	"div-dash/internal/coingecko"
	"div-dash/internal/config"
	"div-dash/internal/controllers"
	"div-dash/internal/iex"
	"div-dash/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	config.ReadConfig()
	config.InitDB()
	r := gin.Default()
	r.Use(gin.Recovery())

	services.JobService().RunJob(iex.IEXExchangesImportJob, services.IEXService().SaveExchanges)
	services.JobService().RunJob(iex.IEXImportSymbolsJob, services.IEXService().SaveSymbols)
	services.JobService().RunJob(iex.ISINAndWKNImportJob, services.IEXService().ImportISINAndWKN)
	services.JobService().RunJob(coingecko.CoingeckoImportCoinsJob, services.CoingeckoService().ImportCryptoSymbols)

	controllers.RegisterRoutes(r)
	port := viper.GetString("server.port")
	r.Run(":" + port)
}
