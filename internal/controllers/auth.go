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
		AbortBadRequest(c, err.Error())
		return
	}

	user, err := config.Queries().FindByEmail(c, loginRequest.Email)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			Abort(c, http.StatusUnauthorized, "wrong credentials")
			return
		}
		c.Error(err)
		return
	}

	if user.Status != db.UserStatusActivated {
		Abort(c, http.StatusUnauthorized, "User not activated")
		return
	}

	if !security.VerifyHash(loginRequest.Password, user.PasswordHash) {
		Abort(c, http.StatusUnauthorized, "wrong credentials")
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
		AbortBadRequest(c, err.Error())
		return
	}

	exists, err := config.Queries().ExistsByEmail(c, registerRequest.Email)

	if err != nil {
		c.Error(err)
		return
	}

	if exists {
		Abort(c, http.StatusConflict, "A user with email '"+registerRequest.Email+"' already exists")
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
		return
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

	registration, err := config.Queries().CreateUserRegistration(c, createRegistrationParams)
	if err != nil {
		c.Error(err)
		return
	}

	body := "Please activate your account at localhost:8080/activate?id=" + registration.ID.String()

	err = services.MailService().SendMail(user.Email, "no-reply@div-dash.io", "Activate your account", body)

	if err != nil {
		c.Error(err)
		return
	}
	c.Status(200)
}

func PostActivate(c *gin.Context) {
	id := c.Query("id")
	registerRequest, err := uuid.Parse(id)

	if err != nil {
		AbortBadRequest(c, "Activation id is in wrong format")
		return
	}

	userRegistration, err := config.Queries().GetUserRegistration(c, registerRequest)

	if err != nil {
		AbortBadRequest(c, "Invalid id")
		return
	}

	if userRegistration.Timestamp.Add(24 * time.Hour).Before(time.Now()) {
		AbortBadRequest(c, "Registration expired")
		return
	}

	err = config.Queries().ActivateUser(c, userRegistration.UserID)

	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
