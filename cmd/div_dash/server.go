package main

import (
	"div-dash/internal/config"
	"div-dash/internal/controllers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.ReadConfig()
	config.InitDB()
	r := gin.Default()
	r.Use(gin.Recovery())
	controllers.RegisterRoutes(r)
	port := viper.GetString("server.port")
	r.Run(":" + port)
}
