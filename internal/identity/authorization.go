package identity

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) authRequired(c *gin.Context) {
	tokenCookie, err := c.Request.Cookie("token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "error": err.Error()})
		return
	}

	token := tokenCookie.Value

	result, userId, err := h.TokenService().VerifyToken(token)

	if result && err == nil {
		c.Set("userId", userId)
		c.Next()
		return
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "error": err.Error()})

}
