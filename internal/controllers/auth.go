package controllers

import (
	"div-dash/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func PostLogin(c *gin.Context) {
	var loginRequest LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.UserService().FindByEmail(loginRequest.Email)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "wrong credentials"})
			return
		}
		c.Error(err)
		return
	}

	if user.Password != loginRequest.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "wrong credentials"})
		return
	}

	token, err := services.TokenService().GenerateToken(user.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
