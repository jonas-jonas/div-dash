package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/db"
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

	for _, balance := range balances {
		symbol := balance.Symbol

		costBasis, err := config.Queries().GetCostBasis(c, db.GetCostBasisParams{
			Symbol: symbol,
			UserID: userId,
		})
		if err != nil {
			c.Error(err)
			return
		}
		resp = append(resp, balanceResponse{
			Symbol:    symbol,
			Total:     balance.Total,
			CostBasis: money.New(costBasis, "EUR").AsMajorUnits(),
		})
	}

	c.JSON(http.StatusOK, resp)
}
