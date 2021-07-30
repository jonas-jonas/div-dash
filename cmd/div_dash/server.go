package main

import (
	"div-dash/internal/config"
	"div-dash/internal/controllers"
	"div-dash/internal/iex"
	"div-dash/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.ReadConfig()
	config.InitDB()
	r := gin.Default()
	r.Use(gin.Recovery())

	services.JobService().RunJob(iex.IEXImportSymbolsJob, services.IEXService().SaveSymbols)
	services.JobService().RunJob(iex.IEXExchangesImportJob, services.IEXService().SaveExchanges)

	controllers.RegisterRoutes(r)
	port := viper.GetString("server.port")
	r.Run(":" + port)
}
