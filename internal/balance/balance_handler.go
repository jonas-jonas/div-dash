package balance

import (
	"div-dash/internal/db"
	"log"
	"net/http"

	"github.com/Rhymond/go-money"
	"github.com/gin-gonic/gin"
)

type (
	balanceHandlerDependencies interface {
		db.QueriesProvider
	}

	BalanceHandlerProvider interface {
		BalanceHandler() *BalanceHandler
	}

	BalanceHandler struct {
		balanceHandlerDependencies
	}
)

func NewBalanceHandler(b balanceHandlerDependencies) *BalanceHandler {
	return &BalanceHandler{
		balanceHandlerDependencies: b,
	}
}

func (b *BalanceHandler) RegisterProtectedRoutes(api gin.IRoutes) {
	api.GET("/balance", b.getBalance)
}

type symbolResponse struct {
	SymbolID   string `json:"symbolID"`
	Type       string `json:"type"`
	Source     string `json:"source"`
	Precision  int32  `json:"precision"`
	SymbolName string `json:"symbolName"`
	ISIN       string `json:"isin"`
	WKN        string `json:"wkn"`
}

func makeSymbolResponse(symbol db.Symbol) symbolResponse {
	return symbolResponse{
		SymbolID:   symbol.SymbolID,
		Type:       symbol.Type,
		Source:     symbol.Source,
		Precision:  symbol.Precision,
		SymbolName: symbol.SymbolName.String,
		ISIN:       symbol.Isin.String,
		WKN:        symbol.Wkn.String,
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

func (b *BalanceHandler) getBalance(c *gin.Context) {
	userId := c.GetString("userId")

	balances, err := b.Queries().GetBalanceByUser(c, userId)
	if err != nil {
		c.Error(err)
		return
	}

	resp := balanceResponse{}

	for _, entry := range balances {
		costBasis := money.New(int64(entry.CostBasis/entry.Amount), "EUR").AsMajorUnits()
		symbol, err := b.Queries().GetSymbol(c, entry.Symbol)
		if err != nil {
			log.Printf("Could not find asset for symbol %s: %s. Skipping balance entry... ", entry.Symbol, err.Error())
			continue
		}
		currentPrice, err := 0.0, nil
		if err != nil {
			log.Printf("Could not get current price for asset %s: %s.", entry.Symbol, err.Error())
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
