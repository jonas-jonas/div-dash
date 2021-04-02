package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/internal/services"
	"div-dash/util/security"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	user, err := config.Queries().FindByEmail(c, loginRequest.Email)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "wrong credentials"})
			return
		}
		c.Error(err)
		return
	}

	if user.Status != db.UserStatusActivated {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not activated"})
		return
	}

	if !security.VerifyHash(loginRequest.Password, user.PasswordHash) {
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

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func PostRegister(c *gin.Context) {
	var registerRequest RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := config.Queries().ExistsByEmail(c, registerRequest.Email)

	if err != nil {
		c.Error(err)
		return
	}

	if exists {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "A user with email '" + registerRequest.Email + "' already exists"})
		return
	}

	passwordHash, err := security.HashPassword(registerRequest.Password)

	if err != nil {
		c.Error(err)
		return
	}

	registerRequestId, err := uuid.NewRandom()

	if err != nil {
		c.Error(err)
	}

	user, err := config.Queries().CreateUser(c, db.CreateUserParams{
		Email:        registerRequest.Email,
		PasswordHash: passwordHash,
		Status:       db.UserStatusRegistered,
	})

	if err != nil {
		c.Error(err)
		return
	}

	createRegistrationParams := db.CreateUserRegistrationParams{
		ID:        registerRequestId,
		UserID:    user.ID,
		Timestamp: time.Now(),
	}

	_, err = config.Queries().CreateUserRegistration(c, createRegistrationParams)
	if err != nil {
		c.Error(err)
		return
	}
	c.Status(200)
}
