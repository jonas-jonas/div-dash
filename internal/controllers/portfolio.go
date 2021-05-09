package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPortfolio(c *gin.Context) {
	id := c.Param("portfolioId")

	portfolio, err := config.Queries().GetPortfolio(c, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			AbortNotFound(c)
		} else {
			c.Error(err)
		}
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
		AbortBadRequest(c, err.Error())
		return
	}

	portfolio, err := config.Queries().CreatePortfolio(c, db.CreatePortfolioParams{
		ID:     "P" + services.IdService().NewId(4),
		Name:   createPortfolioRequest.Name,
		UserID: c.GetString("userId"),
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
	id := c.Param("portfolioId")

	var updatePortfolioRequest updatePortfolioRequest
	if err := c.ShouldBindJSON(&updatePortfolioRequest); err != nil {
		AbortBadRequest(c, err.Error())
		return
	}

	portfolio, err := config.Queries().UpdatePortfolio(c, db.UpdatePortfolioParams{
		ID:   id,
		Name: updatePortfolioRequest.Name,
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

func DeletePortfolio(c *gin.Context) {
	id := c.Param("portfolioId")

	err := config.Queries().DeletePortfolio(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func GetPortfolios(c *gin.Context) {
	userId := c.GetString("userId")

	portfolios, err := config.Queries().ListPortfolios(c, userId)

	if err != nil {
		c.Error(err)
		return
	}

	if portfolios == nil {
		portfolios = []db.Portfolio{}
	}
	c.JSON(http.StatusOK, portfolios)

}
