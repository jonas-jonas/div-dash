package middleware

import (
	"div-dash/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthRequired(t *testing.T) {

	authRequired := AuthRequired()

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, r := gin.CreateTestContext(w)
	token, _ := services.TokenService().GenerateToken(99)

	r.Use(authRequired)
	r.GET("/test", func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			t.Fail()
		}
		assert.Equal(t, 99, userId)
		c.Status(200)
	})

	context.Request, _ = http.NewRequest("GET", "/test", nil)
	context.Request.Header.Add("Authorization", "Bearer "+token)

	r.ServeHTTP(w, context.Request)

	assert.Equal(t, 200, w.Code)
}

func TestAuthRequiredWithInvalidToken(t *testing.T) {

	authRequired := AuthRequired()

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, r := gin.CreateTestContext(w)

	r.Use(authRequired)
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	context.Request, _ = http.NewRequest("GET", "/test", nil)
	context.Request.Header.Add("Authorization", "invalid-token")

	r.ServeHTTP(w, context.Request)

	assert.Equal(t, 401, w.Code)
	assert.JSONEq(t, `{"message": "Unauthorized"}`, w.Body.String())
}
