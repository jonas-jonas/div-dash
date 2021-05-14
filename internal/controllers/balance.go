package controllers

import (
	"div-dash/internal/config"
	"net/http"

	"github.com/Rhymond/go-money"
	"github.com/gin-gonic/gin"
)

type balanceResponse struct {
	Symbol    string  `json:"symbol"`
	Total     float64 `json:"total"`
	CostBasis float64 `json:"costBasis"`
}

func GetBalance(c *gin.Context) {
	userId := c.GetString("userId")

	balances, err := config.Queries().GetBalance(c, userId)
	if err != nil {
		c.Error(err)
		return
	}

	resp := []balanceResponse{}

	for _, entry := range balances {
		costBasis := entry.CostBasis / entry.Total
		resp = append(resp, balanceResponse{
			Symbol:    entry.Symbol,
			Total:     entry.Total,
			CostBasis: money.New(int64(costBasis), "EUR").AsMajorUnits(),
		})
	}

	c.JSON(http.StatusOK, resp)
}
