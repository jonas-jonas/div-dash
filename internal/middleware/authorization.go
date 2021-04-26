package middleware

import (
	"div-dash/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "error": err.Error()})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		result, userId, err := services.TokenService().VerifyToken(token)

		if result && err == nil {
			c.Set("userId", userId)
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "error": err.Error()})
	}
}
