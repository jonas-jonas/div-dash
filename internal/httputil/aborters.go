package httputil

import (
	"div-dash/internal/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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

	c.AbortWithStatusJSON(status, model.APIError{
		Timestamp: timestamp,
		Status:    status,
		Message:   message,
		Path:      path,
	})
}
