package main

import (
	"div-dash/internal/config"
	"div-dash/internal/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	r := gin.Default()
	r.Use(gin.Recovery())
	controllers.RegisterRoutes(r)
	r.Run(":8080")
}
