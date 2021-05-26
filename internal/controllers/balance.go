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
	PLAbsolute     float64  `json:"plAbsolute"`
	PLPercent      float64  `json:"plPercent"`
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
		costBasis := money.New(int64(entry.CostBasis), "EUR").AsMajorUnits()
		asset, err := config.Queries().GetAsset(c, entry.Symbol)
		if err != nil {
			config.Logger().Printf("Could not find asset for symbol %s: %s. Skipping balance entry... ", entry.Symbol, err.Error())
			continue
		}
		currentPrice, err := services.PriceService().GetPriceOfAsset(asset)
		if err != nil {
			config.Logger().Printf("Could not get current price for asset %s: %s.", entry.Symbol, err.Error())
			currentPrice = -0.0
		}
		fiatValue := currentPrice * entry.Amount
		plAbsolute := fiatValue - costBasis
		plPercent := plAbsolute / costBasis
		resp = append(resp, balanceResponse{
			Asset:          asset,
			Amount:         entry.Amount,
			CostBasis:      costBasis,
			FiatAssetPrice: currentPrice,
			FiatValue:      fiatValue,
			PLAbsolute:     plAbsolute,
			PLPercent:      plPercent,
		})
	}

	c.JSON(http.StatusOK, resp)
}
