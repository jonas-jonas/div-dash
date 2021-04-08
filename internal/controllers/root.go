package controllers

import (
	"div-dash/internal/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

//
// Middleware Error Handler in server package
//
func handleErrors(c *gin.Context) {
	c.Next() // execute all the handlers

	// at this point, all the handlers finished. Let's read the errors!
	// in this example we only will use the **last error typed as public**
	// but you could iterate over all them since c.Errors is a slice!
	errorToPrint := c.Errors.Last()
	if errorToPrint != nil {
		log.Printf("Caught error on %s: %s", c.Request.RequestURI, errorToPrint.Error())
		c.JSON(500, gin.H{
			"status":  500,
			"message": errorToPrint.Error(),
		})
	}
}
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(handleErrors)
	api.GET("/ping", Ping)
	api.POST("/login", PostLogin)
	api.POST("/register", PostRegister)
	api.GET("/activate", PostActivate)

	authorized := api.Group("/")
	authorized.Use(middleware.AuthRequired())
	{
		authorized.GET("/user/:id", GetUser)
		authorized.POST("/user/", PostUser)

		authorized.POST("/portfolio", PostPortfolio)
		authorized.GET("/portfolio", GetPortfolios)
		authorized.GET("/portfolio/:id", GetPortfolio)
		authorized.PUT("/portfolio/:id", PutPortfolio)
		authorized.DELETE("/portfolio/:id", DeletePortfolio)
	}
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
