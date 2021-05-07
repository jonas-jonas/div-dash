package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func userResponseFromUser(user db.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}

func GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := config.Queries().GetUser(c, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			AbortNotFound(c)
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, userResponseFromUser(user))
}
