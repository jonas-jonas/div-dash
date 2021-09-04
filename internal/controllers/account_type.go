package controllers

import (
	"database/sql"
	"div-dash/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type accountTypeResponse struct {
	AccountType string `json:"accountType"`
	Label       string `json:"label"`
}

func GetAccountTypes(c *gin.Context) {

	accountTypes, err := config.Queries().ListAccountTypes(c)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, []accountTypeResponse{})
		}
	}

	c.JSON(http.StatusOK, accountTypes)
}
