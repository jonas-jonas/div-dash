package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/internal/services"
	"net/http"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type transactionResponse struct {
	ID                  string          `json:"transactionId"`
	Symbol              string          `json:"symbol"`
	Type                string          `json:"type"`
	TransactionProvider string          `json:"transactionProvider"`
	BuyIn               float64         `json:"buyIn"`
	BuyInDate           time.Time       `json:"buyInDate"`
	Amount              decimal.Decimal `json:"amount"`
	PortfolioId         string          `json:"portfolioId"`
	Side                string          `json:"side"`
}

func marshalTransactionResponse(transaction db.Transaction) transactionResponse {
	return transactionResponse{
		ID:                  transaction.ID,
		Symbol:              transaction.Symbol,
		Type:                string(transaction.Type),
		TransactionProvider: string(transaction.TransactionProvider),
		BuyIn:               money.New(transaction.BuyIn, "EUR").AsMajorUnits(),
		BuyInDate:           transaction.BuyInDate,
		Amount:              transaction.Amount,
		PortfolioId:         transaction.PortfolioID,
		Side:                transaction.Side,
	}
}

func GetTransaction(c *gin.Context) {
	// TODO: Check permissions
	transactionId := c.Param("transactionId")

	transaction, err := config.Queries().GetTransaction(c, transactionId)

	if err != nil {
		c.Error(err)
		return
	}

	resp := marshalTransactionResponse(transaction)

	c.JSON(http.StatusOK, resp)
}

type createTransactionRequest struct {
	Symbol              string    `json:"symbol" binding:"required"`
	Type                string    `json:"type" binding:"required"`
	TransactionProvider string    `json:"transactionProvider" binding:"required"`
	BuyIn               float64   `json:"buyIn" binding:"required"`
	BuyInDate           time.Time `json:"buyInDate" binding:"required"`
	Amount              float64   `json:"amount" binding:"required"`
	Side                string    `json:"side" binding:"required"`
}

func PostTransaction(c *gin.Context) {
	// TODO: Check permissions

	portfolioId := c.Param("portfolioId")

	var createTransactionRequest createTransactionRequest
	if err := c.ShouldBindJSON(&createTransactionRequest); err != nil {
		AbortBadRequest(c, err.Error())
		return
	}

	params := db.CreateTransactionParams{
		ID:                  "T" + services.IdService().NewId(5),
		Symbol:              createTransactionRequest.Symbol,
		Type:                createTransactionRequest.Type,
		TransactionProvider: createTransactionRequest.TransactionProvider,
		BuyIn:               int64(createTransactionRequest.BuyIn * 100),
		BuyInDate:           createTransactionRequest.BuyInDate,
		Amount:              decimal.NewFromFloat(createTransactionRequest.Amount),
		PortfolioID:         portfolioId,
		Side:                createTransactionRequest.Side,
	}

	transaction, err := config.Queries().CreateTransaction(c, params)
	if err != nil {
		c.Error(err)
		return
	}

	resp := marshalTransactionResponse(transaction)

	c.JSON(http.StatusOK, resp)
}

func GetTransactions(c *gin.Context) {
	// TODO: Check permissions

	portfolioId := c.Param("portfolioId")

	transactions, err := config.Queries().ListTransactions(c, portfolioId)

	if err != nil {
		c.Error(err)
		return
	}

	result := []transactionResponse{}

	for _, transaction := range transactions {
		result = append(result, marshalTransactionResponse(transaction))
	}

	c.JSON(http.StatusOK, result)
}
