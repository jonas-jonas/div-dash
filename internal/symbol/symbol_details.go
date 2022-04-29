package symbol

import (
	"context"
	"div-dash/internal/httputil"
	"div-dash/internal/model"

	"github.com/gin-gonic/gin"
)

func (s *SymbolHandler) getSymbolDetails(c *gin.Context) {

	symbolId := c.Param("symbolId")
	symbol, err := s.Queries().GetSymbol(c, symbolId)

	if err != nil {
		httputil.AbortServerError(c)
		return
	}

	var details model.SymbolDetails

	if symbol.Type == "crypto" {
		details, err = s.CoingeckoService().GetDetails(c, symbol)
		if err != nil {
			httputil.AbortBadRequest(c, "Could not get details from coingecko of "+symbolId)
			return
		}
	} else {
		details, err = s.IEXService().GetDetails(context.Background(), symbol)
		if err != nil {
			httputil.AbortBadRequest(c, "Could not get details from iex of "+symbolId)
			return
		}
	}

	c.JSON(200, details)
}
