package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func userResponseFromUser(user db.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}

func GetUser(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		AbortBadRequest(c, "User id is invalid")
		return
	}

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
