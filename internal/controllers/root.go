package controllers

import (
	"div-dash/internal/middleware"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Message   string    `json:"message"`
	Status    int       `json:"status"`
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
}

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
		AbortServerError(c)
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

func AbortBadRequest(c *gin.Context, message string) {
	Abort(c, http.StatusBadRequest, message)
}

func AbortServerError(c *gin.Context) {
	Abort(c, http.StatusInternalServerError, "An internal server error occured. Please try again later.")
}

func AbortUnauthorized(c *gin.Context) {
	Abort(c, http.StatusUnauthorized, "Please log in and try again")
}

func AbortNotFound(c *gin.Context) {
	Abort(c, http.StatusNotFound, "The requested resource could not be found")
}

func Abort(c *gin.Context, status int, message string) {
	timestamp := time.Now()
	path := c.Request.URL.Path

	c.AbortWithStatusJSON(status, APIError{
		Timestamp: timestamp,
		Status:    status,
		Message:   message,
		Path:      path,
	})
}
