package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/internal/services"
	"net/http"

	"github.com/Rhymond/go-money"
	"github.com/gin-gonic/gin"
)

type balanceResponse struct {
	Asset          db.Asset `json:"asset"`
	Amount         float64  `json:"amount"`
	CostBasis      float64  `json:"costBasis"`
	FiatAssetPrice float64  `json:"fiatAssetPrice"`
	FiatValue      float64  `json:"fiatValue"`
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
		asset, err := config.Queries().GetAsset(c, entry.Symbol)
		if err != nil {
			config.Logger().Printf("Could not find asset for symbol %s: %s. Skipping balance entry... ", entry.Symbol, err.Error())
			continue
		}
		currentPrice, err := services.PriceService().GetPriceOfAsset(c, asset)
		if err != nil {
			config.Logger().Printf("Could not get current price for asset %s: %s.", entry.Symbol, err.Error())
			currentPrice = -0.0
		}
		resp = append(resp, balanceResponse{
			Asset:          asset,
			Amount:         entry.Amount,
			CostBasis:      money.New(int64(costBasis), "EUR").AsMajorUnits(),
			FiatAssetPrice: currentPrice,
			FiatValue:      currentPrice * entry.Amount,
		})
	}

	c.JSON(http.StatusOK, resp)
}
