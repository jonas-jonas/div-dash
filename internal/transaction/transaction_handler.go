package transaction

import (
	"div-dash/internal/csvimport"
	"div-dash/internal/db"
	"div-dash/internal/httputil"
	"div-dash/internal/id"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type (
	transactionHandlerDependencies interface {
		db.QueriesProvider
		id.IdServiceProvider
		csvimport.CSVImporterProvider
	}

	TransactionHandlerProvider interface {
		TransactionHandler() *TransactionHandler
	}

	TransactionHandler struct {
		transactionHandlerDependencies
	}
)

func NewTransactionHandler(t transactionHandlerDependencies) *TransactionHandler {
	return &TransactionHandler{
		transactionHandlerDependencies: t,
	}
}

func (t *TransactionHandler) RegisterProtectedRoutes(api gin.IRoutes) {
	api.POST("/account/:accountId/transaction", t.postTransaction)
	api.GET("/account/:accountId/transaction", t.getTransactions)
	api.GET("/account/:accountId/transaction/:transactionId", t.getTransaction)
	api.POST("/account/:accountId/transaction-import", t.postAccountTransactionImport)
}

type transactionResponse struct {
	ID                  string    `json:"transactionId"`
	Symbol              string    `json:"symbol"`
	Type                string    `json:"type"`
	TransactionProvider string    `json:"transactionProvider"`
	Price               float64   `json:"price"`
	Date                time.Time `json:"date"`
	Amount              float64   `json:"amount"`
	AccountId           string    `json:"accountId"`
	Side                string    `json:"side"`
}

func marshalTransactionResponse(transaction db.Transaction) transactionResponse {
	amount, _ := transaction.Amount.Float64()
	return transactionResponse{
		ID:                  transaction.ID,
		Symbol:              transaction.Symbol,
		Type:                string(transaction.Type),
		TransactionProvider: string(transaction.TransactionProvider),
		Price:               money.New(transaction.Price, "EUR").AsMajorUnits(),
		Date:                transaction.Date,
		Amount:              amount,
		AccountId:           transaction.AccountID,
		Side:                transaction.Side,
	}
}

func (t *TransactionHandler) getTransaction(c *gin.Context) {
	// TODO: Check permissions
	transactionId := c.Param("transactionId")
	accountId := c.Param("accountId")
	userId := c.GetString("userId")

	transaction, err := t.Queries().GetTransaction(c, db.GetTransactionParams{
		ID:        transactionId,
		AccountID: accountId,
		UserID:    userId,
	})

	if err != nil {
		c.Error(err)
		return
	}
	decimal.MarshalJSONWithoutQuotes = true

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

func (t *TransactionHandler) postTransaction(c *gin.Context) {
	// TODO: Check permissions

	accountId := c.Param("accountId")
	userId := c.GetString("userId")

	var createTransactionRequest createTransactionRequest
	if err := c.ShouldBindJSON(&createTransactionRequest); err != nil {
		httputil.AbortBadRequest(c, err.Error())
		return
	}

	params := db.CreateTransactionParams{
		ID:                  "T" + t.IdService().NewID(5),
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

	transactionId, err := t.Queries().CreateTransaction(c, params)
	if err != nil {
		c.Error(err)
		return
	}

	transaction, err := t.Queries().GetTransaction(c, db.GetTransactionParams{
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

func (t *TransactionHandler) getTransactions(c *gin.Context) {

	accountId := c.Param("accountId")
	userId := c.GetString("userId")

	transactions, err := t.Queries().ListTransactions(c, db.ListTransactionsParams{
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

func (h *TransactionHandler) postAccountTransactionImport(c *gin.Context) {
	userId := c.GetString("userId")
	accountId := c.Param("accountId")

	file, headers, err := c.Request.FormFile("file")

	if err != nil {
		httputil.AbortBadRequest(c, "file err "+err.Error())
		return
	}

	if ext := filepath.Ext(headers.Filename); ext != ".xls" {
		httputil.AbortBadRequest(c, "Unsupported file type: "+ext)
		return
	}

	err = h.CSVImporter().ImportCSV(c, file, accountId, userId)

	if err != nil {
		httputil.AbortBadRequest(c, "Could not parse file: "+err.Error())
		return
	}

	transactions, err := h.Queries().ListTransactions(c, db.ListTransactionsParams{
		AccountID: accountId,
		UserID:    userId,
	})

	if err != nil {
		c.Error(err)
		return
	}

	var response []transactionResponse

	for _, transaction := range transactions {
		response = append(response, marshalTransactionResponse(transaction))
	}

	c.JSON(http.StatusOK, response)
}
