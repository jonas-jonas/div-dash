package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPortfolio(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	portfolio, err := config.Queries().GetPortfolio(c, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"message": "Portfolio with id '" + idString + "' not found"})
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

type createPortfolioRequest struct {
	Name string `json:"name" binding:"required"`
}

func PostPortfolio(c *gin.Context) {
	var createPortfolioRequest createPortfolioRequest

	if err := c.ShouldBindJSON(&createPortfolioRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	portfolio, err := config.Queries().CreatePortfolio(c, db.CreatePortfolioParams{
		Name:   createPortfolioRequest.Name,
		UserID: c.GetInt64("userId"),
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

type updatePortfolioRequest struct {
	Name string `json:"name" binding:"required"`
}

func PutPortfolio(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	var updatePortfolioRequest updatePortfolioRequest
	if err := c.ShouldBindJSON(&updatePortfolioRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	portfolio, err := config.Queries().UpdatePortfolio(c, db.UpdatePortfolioParams{
		PortfolioID: id,
		Name:        updatePortfolioRequest.Name,
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

func DeletePortfolio(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	err = config.Queries().DeletePortfolio(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func GetPortfolios(c *gin.Context) {
	userId := c.GetInt64("userId")

	portfolios, err := config.Queries().ListPortfolios(c, userId)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, portfolios)
}
