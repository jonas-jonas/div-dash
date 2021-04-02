package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/util/security"
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

func PostUser(c *gin.Context) {
	var createUserRequest CreateUserRequest

	if err := c.ShouldBindJSON(&createUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := config.Queries().ExistsByEmail(c, createUserRequest.Email)

	if err != nil {
		c.Error(err)
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"message": "A user with email '" + createUserRequest.Email + "' already exists"})
		return
	}

	passwordHash, err := security.HashPassword(createUserRequest.Password)

	if err != nil {
		c.Error(err)
		return
	}

	user, err := config.Queries().CreateUser(c, db.CreateUserParams{
		Email:        createUserRequest.Email,
		PasswordHash: passwordHash,
		Status:       db.UserStatusActivated,
	})

	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, userResponseFromUser(user))
}

func GetUser(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Param 'id' must be int, got '" + idString + "'"})
		return
	}

	user, err := config.Queries().GetUser(c, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"message": "User with id '" + idString + "' not found"})
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, userResponseFromUser(user))
}
