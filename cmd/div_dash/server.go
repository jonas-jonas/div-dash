package main

import (
	"div-dash/internal/config"
	"div-dash/internal/controllers"
	"div-dash/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.ReadConfig()
	config.InitDB()
	r := gin.Default()
	r.Use(gin.Recovery())
	err := services.BinanceService().SaveExchangeInfo()
	if err != nil {
		log.Printf("Error Importing Binance Exchange Info %s", err.Error())
	}

	err = services.IEXService().SaveSymbols()
	if err != nil {
		log.Printf("Error saving IEX Symbols %s", err.Error())
	}

	err = services.IEXService().SaveExchanges()
	if err != nil {
		log.Printf("Error saving IEX Exchanges %s", err.Error())
	}
	controllers.RegisterRoutes(r)
	port := viper.GetString("server.port")
	r.Run(":" + port)
}
