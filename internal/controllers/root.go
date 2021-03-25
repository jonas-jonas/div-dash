package controllers

import (
	"div-dash/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/ping", Ping)
	api.POST("/login", PostLogin)

	authorized := api.Group("/")
	authorized.Use(middleware.AuthRequired())
	{
		authorized.GET("/user/:id", GetUser)
		authorized.POST("/user/", PostUser)
	}
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
