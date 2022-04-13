package account

import (
	"database/sql"
	"div-dash/internal/db"
	"div-dash/internal/httputil"
	"div-dash/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	accountHandlerDependencies interface {
		db.QueriesProvider
	}

	AccountHandlerProvider interface {
		AccountHandler() *Handler
	}

	Handler struct {
		accountHandlerDependencies
	}
)

func NewAccountHandler(a accountHandlerDependencies) *Handler {
	return &Handler{accountHandlerDependencies: a}
}

func (h *Handler) RegisterProtectedRoutes(api gin.IRoutes) {
	api.POST("/account", h.postAccount)
	api.GET("/account", h.getAccounts)
	api.GET("/account/:accountId", h.getAccount)
	api.PUT("/account/:accountId", h.putAccount)
	api.DELETE("/account/:accountId", h.deleteAccount)

	api.GET("/account-types", h.getAccountTypes)
}

func (h *Handler) getAccount(c *gin.Context) {
	id := c.Param("accountId")

	account, err := h.Queries().GetAccount(c, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			httputil.AbortNotFound(c)
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

func (h *Handler) postAccount(c *gin.Context) {
	var createAccountRequest createAccountRequest

	if err := c.ShouldBindJSON(&createAccountRequest); err != nil {
		httputil.AbortBadRequest(c, err.Error())
		return
	}

	account, err := h.Queries().CreateAccount(c, db.CreateAccountParams{
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

func (h *Handler) putAccount(c *gin.Context) {
	id := c.Param("accountId")

	var updateAccountRequest updateAccountRequest
	if err := c.ShouldBindJSON(&updateAccountRequest); err != nil {
		httputil.AbortBadRequest(c, err.Error())
		return
	}

	account, err := h.Queries().UpdateAccount(c, db.UpdateAccountParams{
		ID:   id,
		Name: updateAccountRequest.Name,
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *Handler) deleteAccount(c *gin.Context) {
	id := c.Param("accountId")

	err := h.Queries().DeleteAccount(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) getAccounts(c *gin.Context) {
	userId := c.GetString("userId")

	accounts, err := h.Queries().ListAccounts(c, userId)

	if err != nil {
		c.Error(err)
		return
	}

	if accounts == nil {
		accounts = []db.Account{}
	}
	c.JSON(http.StatusOK, accounts)

}
