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
	Price               float64         `json:"price"`
	Date                time.Time       `json:"date"`
	Amount              decimal.Decimal `json:"amount"`
	AccountId           string          `json:"accountId"`
	Side                string          `json:"side"`
}

func marshalTransactionResponse(transaction db.Transaction) transactionResponse {
	return transactionResponse{
		ID:                  transaction.ID,
		Symbol:              transaction.Symbol,
		Type:                string(transaction.Type),
		TransactionProvider: string(transaction.TransactionProvider),
		Price:               money.New(transaction.Price, "EUR").AsMajorUnits(),
		Date:                transaction.Date,
		Amount:              transaction.Amount,
		AccountId:           transaction.AccountID,
		Side:                transaction.Side,
	}
}

func GetTransaction(c *gin.Context) {
	// TODO: Check permissions
	transactionId := c.Param("transactionId")
	accountId := c.Param("accountId")
	userId := c.GetString("userId")

	transaction, err := config.Queries().GetTransaction(c, db.GetTransactionParams{
		ID:        transactionId,
		AccountID: accountId,
		UserID:    userId,
	})

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
	Price               float64   `json:"price" binding:"required"`
	Date                time.Time `json:"date" binding:"required"`
	Amount              float64   `json:"amount" binding:"required"`
	Side                string    `json:"side" binding:"required"`
}

func PostTransaction(c *gin.Context) {
	// TODO: Check permissions

	accountId := c.Param("accountId")
	userId := c.GetString("userId")

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
		Price:               int64(createTransactionRequest.Price * 100),
		Date:                createTransactionRequest.Date,
		Amount:              decimal.NewFromFloat(createTransactionRequest.Amount),
		AccountID:           accountId,
		UserID:              userId,
		Side:                createTransactionRequest.Side,
	}

	transactionId, err := config.Queries().CreateTransaction(c, params)
	if err != nil {
		c.Error(err)
		return
	}

	transaction, err := config.Queries().GetTransaction(c, db.GetTransactionParams{
		ID:        transactionId,
		AccountID: accountId,
		UserID:    userId,
	})
	if err != nil {
		c.Error(err)
		return
	}

	resp := marshalTransactionResponse(transaction)

	c.JSON(http.StatusOK, resp)
}

func GetTransactions(c *gin.Context) {

	accountId := c.Param("accountId")
	userId := c.GetString("userId")

	transactions, err := config.Queries().ListTransactions(c, db.ListTransactionsParams{
		AccountID: accountId,
		UserID:    userId,
	})

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
