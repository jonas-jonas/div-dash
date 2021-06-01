package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SearchSymbol(c *gin.Context) {
	search := c.Query("query")
	count, err := strconv.ParseInt(c.Query("count"), 10, 32)
	if err != nil {
		AbortBadRequest(c, "param 'count' must be int")
		return
	}

	symbols, err := config.Queries().SearchSymbol(c, db.SearchSymbolParams{
		Search: "%" + search + "%",
		Count:  int32(count),
	})
	if err != nil {
		c.Error(err)
		return
	}

	resp := []symbolResponse{}

	for _, symbol := range symbols {
		resp = append(resp, symbolResponse{
			SymbolID:   symbol.SymbolID,
			Type:       symbol.Type,
			Source:     symbol.Source,
			Precision:  symbol.Precision,
			SymbolName: symbol.SymbolName.String,
		})
	}

	c.JSON(http.StatusOK, resp)
}
