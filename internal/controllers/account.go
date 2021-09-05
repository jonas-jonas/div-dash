package controllers

import (
	"database/sql"
	"div-dash/internal/account_types/comdirect"
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/internal/services"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func GetAccount(c *gin.Context) {
	id := c.Param("accountId")

	account, err := config.Queries().GetAccount(c, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			AbortNotFound(c)
		} else {
			c.Error(err)
		}
		return
	}

	c.JSON(http.StatusOK, account)
}

type createAccountRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"accountType"`
}

func PostAccount(c *gin.Context) {
	var createAccountRequest createAccountRequest

	if err := c.ShouldBindJSON(&createAccountRequest); err != nil {
		AbortBadRequest(c, err.Error())
		return
	}

	account, err := config.Queries().CreateAccount(c, db.CreateAccountParams{
		ID:     "A" + services.IdService().NewId(4),
		Name:   createAccountRequest.Name,
		UserID: c.GetString("userId"),
		AccountType: sql.NullString{
			String: createAccountRequest.Type,
			Valid:  true,
		},
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, account)
}

type updateAccountRequest struct {
	Name string `json:"name" binding:"required"`
}

func PutAccount(c *gin.Context) {
	id := c.Param("accountId")

	var updateAccountRequest updateAccountRequest
	if err := c.ShouldBindJSON(&updateAccountRequest); err != nil {
		AbortBadRequest(c, err.Error())
		return
	}

	account, err := config.Queries().UpdateAccount(c, db.UpdateAccountParams{
		ID:   id,
		Name: updateAccountRequest.Name,
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, account)
}

func DeleteAccount(c *gin.Context) {
	id := c.Param("accountId")

	err := config.Queries().DeleteAccount(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func GetAccounts(c *gin.Context) {
	userId := c.GetString("userId")

	accounts, err := config.Queries().ListAccounts(c, userId)

	if err != nil {
		c.Error(err)
		return
	}

	if accounts == nil {
		accounts = []db.Account{}
	}
	c.JSON(http.StatusOK, accounts)

}

func PostAccountTransactionImport(c *gin.Context) {
	userId := c.GetString("userId")
	accountId := c.Param("accountId")

	file, headers, err := c.Request.FormFile("file")

	if err != nil {
		AbortBadRequest(c, "file err "+err.Error())
		return
	}

	if ext := filepath.Ext(headers.Filename); ext != ".xls" {
		AbortBadRequest(c, "Unsupported file type: "+ext)
		return
	}

	importer := comdirect.NewCsvImporter(config.Queries(), config.DB())

	err = importer.ImportTransactionsXLS(c, file, accountId, userId)

	if err != nil {
		AbortBadRequest(c, "Could not parse file: "+err.Error())
		return
	}

	transactions, err := config.Queries().ListTransactions(c, db.ListTransactionsParams{
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
