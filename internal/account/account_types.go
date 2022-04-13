package account

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type accountTypeResponse struct {
	AccountType string `json:"accountType"`
	Label       string `json:"label"`
}

func (h *Handler) getAccountTypes(c *gin.Context) {

	accountTypes, err := h.Queries().ListAccountTypes(c)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, []accountTypeResponse{})
		}
	}

	c.JSON(http.StatusOK, accountTypes)
}
