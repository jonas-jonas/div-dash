package controllers

import (
	"div-dash/internal/db"
	"div-dash/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func PostUser(c *gin.Context) {
	var createUserRequest CreateUserRequest

	if err := c.ShouldBindJSON(&createUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := services.UserService().ExistsByEmail(createUserRequest.Email)

	if err != nil {
		c.Error(err)
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"message": "A user with email '" + createUserRequest.Email + "' already exists"})
		return
	}

	user, err := services.UserService().CreateUser(db.CreateUserParams(createUserRequest))

	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, UserResponse(user))
}

func GetUser(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Param 'id' must be int, got '" + idString + "'"})
		return
	}

	user, err := services.UserService().FindById(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"message": "User with id '" + idString + "' not found"})
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, UserResponse(user))
}
