package middleware

import (
	"div-dash/internal/services"
	"div-dash/util/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthRequired(t *testing.T) {

	testutil.SetupConfig()
	authRequired := AuthRequired()

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, r := gin.CreateTestContext(w)
	token, _ := services.TokenService().GenerateToken(testutil.TestUserID)

	r.Use(authRequired)
	r.GET("/test", func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			t.Fail()
		}
		assert.Equal(t, testutil.TestUserID, userId)
		c.Status(200)
	})

	context.Request, _ = http.NewRequest("GET", "/test", nil)
	authCookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
	}
	context.Request.AddCookie(&authCookie)

	r.ServeHTTP(w, context.Request)

	assert.Equal(t, 200, w.Code)
}

func TestAuthRequiredWithInvalidToken(t *testing.T) {
	testutil.SetupConfig()
	authRequired := AuthRequired()

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, r := gin.CreateTestContext(w)

	r.Use(authRequired)
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	context.Request, _ = http.NewRequest("GET", "/test", nil)
	authCookie := http.Cookie{
		Name:     "token",
		Value:    "invalid-token",
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
	}
	context.Request.AddCookie(&authCookie)

	r.ServeHTTP(w, context.Request)

	assert.Equal(t, 401, w.Code)
	assert.JSONEq(t, `{"error":"failed to decode token: incorrect token header", "message":"Unauthorized"}`, w.Body.String())
}
