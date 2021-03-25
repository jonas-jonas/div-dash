package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	password string
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var allUsers = []User{
	{
		ID:       0,
		Email:    "test@test.de",
		password: "password",
	},
}

func PostUser(c *gin.Context) {
	var createUserRequest CreateUserRequest

	if err := c.ShouldBindJSON(&createUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, user := range allUsers {
		if user.Email == createUserRequest.Email {
			c.JSON(http.StatusConflict, gin.H{"message": "A user with email '" + createUserRequest.Email + "' already exists"})
			return
		}
	}

	user := User{
		ID:       1,
		Email:    createUserRequest.Email,
		password: createUserRequest.Password,
	}

	allUsers = append(allUsers, user)

	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Param 'id' must be int, got '" + idString + "'"})
		return
	}
	for _, user := range allUsers {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "User with id '" + idString + "' not found"})
}
