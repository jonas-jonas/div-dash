package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/internal/services"
	"div-dash/util/security"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func GetAuthForm(c *gin.Context) {
	ref := c.Request.Referer()
	c.HTML(http.StatusOK, "login.html", gin.H{"data": gin.H{"Referer": ref}})
}

type LoginFormRequest struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
	Referer  string `form:"referer"`
}

func AbortForm(c *gin.Context, field, message string, data interface{}) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"errors": gin.H{
			field: message,
		},
		"data": data,
	})
}

func PostAuthForm(c *gin.Context) {
	var authForm LoginFormRequest
	if err := c.ShouldBind(&authForm); err != nil {
		ve, _ := err.(validator.ValidationErrors)
		for _, e := range ve {
			if e.Tag() == "required" {
				AbortForm(c, e.Field(), e.Field()+" is required", authForm)
				return
			} else {
				AbortForm(c, e.Field(), e.Error(), authForm)
				return
			}
		}
	}
	user, err := config.Queries().FindByEmail(c, authForm.Email)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			AbortForm(c, "Password", "Invalid password", authForm)
			return
		}
		config.Logger().Printf("Failed to fetch user account: %s", err.Error())
		AbortForm(c, "Email", "Internal Server Error.", authForm)
		return
	}
	if user.Status != db.UserStatusActivated {
		AbortForm(c, "Email", "User not activated", authForm)
		return
	}

	if !security.VerifyHash(authForm.Password, user.PasswordHash) {
		AbortForm(c, "Password", "Invalid password", authForm)
		return
	}

	token, err := services.TokenService().GenerateToken(user.ID)
	if err != nil {
		AbortForm(c, "Password", "Internal server error. Please try again later.", nil)
		return
	}
	c.SetCookie("token", token, 24*60*60, "/", "localhost", true, true)
	if authForm.Referer != "" {
		c.Redirect(http.StatusSeeOther, authForm.Referer)
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
}

func GetRegisterForm(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

type RegisterFormRequest struct {
	Email          string `form:"email" binding:"required"`
	Password       string `form:"password" binding:"required"`
	RepeatPassword string `form:"repeatPassword" binding:"required"`
	AcceptTOS      bool   `form:"acceptTOS"`
}

func PostRegisterForm(c *gin.Context) {
	var registerFormRequest RegisterFormRequest

	if err := c.ShouldBind(&registerFormRequest); err != nil {
		ve, _ := err.(validator.ValidationErrors)
		errors := gin.H{}
		for _, e := range ve {
			field := e.Field()
			var message string
			if e.Tag() == "required" {
				message = field + " is required"
			} else {
				message = e.Error()
			}
			errors[field] = message
		}
		c.HTML(http.StatusOK, "register.html", gin.H{
			"errors": errors,
			"data":   registerFormRequest,
		})
		log.Printf("err: %v", err)
		return
	}

	if !registerFormRequest.AcceptTOS {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"errors": gin.H{
				"AcceptTOS": "Must be accepted",
			},
			"data": registerFormRequest,
		})
		return
	}

	if registerFormRequest.Password != registerFormRequest.RepeatPassword {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"errors": gin.H{
				"RepeatPassword": "Must match password",
			},
			"data": registerFormRequest,
		})
		return
	}

	exists, err := config.Queries().ExistsByEmail(c, registerFormRequest.Email)

	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"errors": gin.H{
				"AcceptTOS": "Internal Server Error",
			},
			"data": registerFormRequest,
		})
		return
	}

	if exists {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"errors": gin.H{
				"Email": "A user with email '" + registerFormRequest.Email + "' already exists",
			},
			"data": registerFormRequest,
		})
		return
	}

	passwordHash, err := security.HashPassword(registerFormRequest.Password)

	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"errors": gin.H{
				"AcceptTOS": "Internal Server Error",
			},
			"data": registerFormRequest,
		})
		return
	}

	registerRequestId, err := uuid.NewRandom()

	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"errors": gin.H{
				"AcceptTOS": "Internal Server Error",
			},
			"data": registerFormRequest,
		})
		return
	}

	user, err := config.Queries().CreateUser(c, db.CreateUserParams{
		Email:        registerFormRequest.Email,
		PasswordHash: passwordHash,
		Status:       db.UserStatusRegistered,
	})

	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"errors": gin.H{
				"AcceptTOS": "Internal Server Error",
			},
			"data": registerFormRequest,
		})
		return
	}

	createRegistrationParams := db.CreateUserRegistrationParams{
		ID:        registerRequestId,
		UserID:    user.ID,
		Timestamp: time.Now(),
	}

	registration, err := config.Queries().CreateUserRegistration(c, createRegistrationParams)
	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"errors": gin.H{
				"AcceptTOS": "Internal Server Error",
			},
			"data": registerFormRequest,
		})
		return
	}

	body := "Please activate your account at localhost:8080/activate?id=" + registration.ID.String()

	err = services.MailService().SendMail(user.Email, "no-reply@div-dash.io", "Activate your account", body)

	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"errors": gin.H{
				"AcceptTOS": "Internal Server Error",
			},
			"data": registerFormRequest,
		})
		return
	}
	c.HTML(http.StatusOK, "register.html", gin.H{
		"success": true,
	})
}

func GetActivateForm(c *gin.Context) {
	id := c.Query("id")
	registerRequest, err := uuid.Parse(id)

	if err != nil {
		c.HTML(http.StatusOK, "activate.html", gin.H{
			"status":  "error",
			"message": "Invalid activation id",
		})
		return
	}

	userRegistration, err := config.Queries().GetUserRegistration(c, registerRequest)

	if err != nil {
		c.HTML(http.StatusOK, "activate.html", gin.H{
			"status":  "error",
			"message": "Invalid activation id",
		})
		return
	}

	if userRegistration.Timestamp.Add(24 * time.Hour).Before(time.Now()) {
		c.HTML(http.StatusOK, "activate.html", gin.H{
			"status":  "error",
			"message": "This activation Id expired. Please register again.",
		})
		return
	}

	activated, err := config.Queries().IsUserActivated(c, userRegistration.UserID)

	if err != nil {
		c.HTML(http.StatusOK, "activate.html", gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return
	}

	if activated {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	err = config.Queries().ActivateUser(c, userRegistration.UserID)

	if err != nil {
		c.HTML(http.StatusOK, "activate.html", gin.H{
			"status":  "error",
			"message": "Internal Server Error",
		})
		return
	}

	c.HTML(http.StatusOK, "activate.html", gin.H{
		"status": "success",
	})
}
