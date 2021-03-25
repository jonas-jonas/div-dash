package main

import (
	"div-dash/internal/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())
	controllers.RegisterRoutes(r)
	r.Run(":8080")
}
