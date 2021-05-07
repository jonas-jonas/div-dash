package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/middleware"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/static"
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
	authForm := r.Group("/")

	r.Use(static.Serve("/", static.LocalFile("web/build", true)))
	r.NoRoute(func(c *gin.Context) {
		c.File("web/build/index.html")
	})

	authForm.GET("/login", GetAuthForm)
	authForm.POST("/login", PostAuthForm)
	authForm.GET("/register", GetRegisterForm)
	authForm.POST("/register", PostRegisterForm)
	authForm.GET("/activate", GetActivateForm)

	api := r.Group("/api")
	api.Use(handleErrors)
	api.GET("/ping", Ping)

	authorized := api.Group("/")
	authorized.Use(middleware.AuthRequired())
	{
		authorized.GET("/auth/identity", GetAuthIdentity)
		authorized.GET("/auth/logout", GetLogout)
		authorized.GET("/user/:id", GetUser)

		authorized.POST("/portfolio", PostPortfolio)
		authorized.GET("/portfolio", GetPortfolios)
		authorized.GET("/portfolio/:portfolioId", GetPortfolio)
		authorized.PUT("/portfolio/:portfolioId", PutPortfolio)
		authorized.DELETE("/portfolio/:portfolioId", DeletePortfolio)

		authorized.POST("/portfolio/:portfolioId/transaction", PostTransaction)
		authorized.GET("/portfolio/:portfolioId/transaction", GetTransactions)
		authorized.GET("/portfolio/:portfolioId/transaction/:transactionId", GetTransaction)
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

func GetAuthIdentity(c *gin.Context) {
	userId := c.GetString("userId")

	user, err := config.Queries().GetUser(c, userId)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, userResponseFromUser(user))
}

func GetLogout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", true, true)
}
