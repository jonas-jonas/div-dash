package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/internal/services"
	"net/http"

	"github.com/Rhymond/go-money"
	"github.com/gin-gonic/gin"
)

type symbolResponse struct {
	SymbolID   string `json:"symbolID"`
	Type       string `json:"type"`
	Source     string `json:"source"`
	Precision  int32  `json:"precision"`
	SymbolName string `json:"symbolName"`
}

func makeSymbolResponse(symbol db.Symbol) symbolResponse {
	return symbolResponse{
		SymbolID:   symbol.SymbolID,
		Type:       symbol.Type,
		Source:     symbol.Source,
		Precision:  symbol.Precision,
		SymbolName: symbol.SymbolName.String,
	}
}

type balanceItemResponse struct {
	Symbol         symbolResponse `json:"symbol"`
	Amount         float64        `json:"amount"`
	CostBasis      float64        `json:"costBasis"`
	FiatAssetPrice float64        `json:"fiatAssetPrice"`
	PNL            pnlResponse    `json:"pnl"`
}

type pnlResponse struct {
	PNL        float64 `json:"pnl"`
	PNLPercent float64 `json:"pnlPercent"`
}

type balanceResponse struct {
	Symbols   []balanceItemResponse `json:"symbols"`
	FiatValue float64               `json:"fiatValue"`
	CostBasis float64               `json:"costBasis"`
	PNL       pnlResponse           `json:"pnl"`
}

func GetBalance(c *gin.Context) {
	userId := c.GetString("userId")

	balances, err := config.Queries().GetBalanceByUser(c, userId)
	if err != nil {
		c.Error(err)
		return
	}

	resp := balanceResponse{}

	for _, entry := range balances {
		costBasis := money.New(int64(entry.CostBasis/entry.Amount), "EUR").AsMajorUnits()
		symbol, err := config.Queries().GetSymbol(c, entry.Symbol)
		if err != nil {
			config.Logger().Printf("Could not find asset for symbol %s: %s. Skipping balance entry... ", entry.Symbol, err.Error())
			continue
		}
		currentPrice, err := services.PriceService().GetPriceOfAsset(symbol)
		if err != nil {
			config.Logger().Printf("Could not get current price for asset %s: %s.", entry.Symbol, err.Error())
			currentPrice = 0.0
		}

		currentFiatValue := currentPrice * entry.Amount
		plAbsolute := currentPrice - costBasis
		plPercent := plAbsolute / costBasis

		resp.CostBasis += costBasis * entry.Amount
		resp.FiatValue += currentFiatValue

		resp.Symbols = append(resp.Symbols, balanceItemResponse{
			Symbol:         makeSymbolResponse(symbol),
			Amount:         entry.Amount,
			CostBasis:      costBasis,
			FiatAssetPrice: currentPrice,
			PNL: pnlResponse{
				PNL:        plAbsolute * entry.Amount,
				PNLPercent: plPercent,
			},
		})
	}

	c.JSON(http.StatusOK, resp)
}
