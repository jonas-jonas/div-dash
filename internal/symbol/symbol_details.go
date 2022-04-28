package symbol

import (
	"context"
	"div-dash/internal/httputil"

	"github.com/gin-gonic/gin"
)

func (s *SymbolHandler) getSymbolDetails(c *gin.Context) {

	symbolId := c.Param("symbolId")
	symbol, err := s.Queries().GetSymbol(c, symbolId)

	if err != nil {
		httputil.AbortServerError(c)
		return
	}

	if symbol.Type == "crypto" {
		httputil.AbortBadRequest(c, "Crypto not supported currently")
		return
	}

	details, err := s.IEXService().GetDetails(context.Background(), symbol)

	if err != nil {
		httputil.AbortBadRequest(c, "Could not get details of "+symbolId)
		return
	}

	c.JSON(200, details)
}
