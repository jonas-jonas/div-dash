package controllers

import (
	"div-dash/internal/config"
	"net/http"

	"github.com/Rhymond/go-money"
	"github.com/gin-gonic/gin"
)

type balanceResponse struct {
	Symbol    string  `json:"symbol"`
	Amount    float64 `json:"amount"`
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
		costBasis := entry.CostBasis / entry.Amount
		resp = append(resp, balanceResponse{
			Symbol:    entry.Symbol,
			Amount:    entry.Amount,
			CostBasis: money.New(int64(costBasis), "EUR").AsMajorUnits(),
		})
	}

	c.JSON(http.StatusOK, resp)
}
