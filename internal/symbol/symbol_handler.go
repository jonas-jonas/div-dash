package symbol

import (
	"div-dash/internal/db"
	"div-dash/internal/httputil"
	"div-dash/internal/iex"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	symbolHandlerDependencies interface {
		db.QueriesProvider
		iex.IEXServiceProvider
	}

	SymbolHandlerProvider interface {
		SymbolHandler() *SymbolHandler
	}

	SymbolHandler struct {
		symbolHandlerDependencies
	}
)

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

func NewSymbolHandler(s symbolHandlerDependencies) *SymbolHandler {
	return &SymbolHandler{
		symbolHandlerDependencies: s,
	}
}

func (s *SymbolHandler) RegisterProtectedRoutes(api gin.IRoutes) {
	api.GET("/symbols", s.getSymbols)
	api.GET("/symbol/search", s.searchSymbol)
	api.GET("/symbol/details/:symbolId", s.getSymbolDetails)
	// api.GET("/symbol/chart/:symbolId", s.getSymbolChart)
}

type paginatedSymbols struct {
	Symbols    []symbolResponse `json:"symbols"`
	TotalCount int64            `json:"totalCount"`
	Pages      int              `json:"pages"`
	ActivePage int              `json:"activePage"`
}

func (s *SymbolHandler) getSymbols(c *gin.Context) {
	symbolType := c.Query("type")

	count, err := strconv.ParseInt(c.Query("count"), 10, 32)
	if err != nil {
		httputil.AbortBadRequest(c, "param 'count' must be int")
		return
	}

	var symbols []db.Symbol
	var totalCount int64

	if symbolType == "" {
		symbols, err = s.Queries().GetSymbols(c, int32(count))
		if err != nil {
			c.Error(err)
			return
		}
		totalCount, err = s.Queries().GetSymbolCount(c)
		if err != nil {
			c.Error(err)
			return
		}
	} else {
		symbols, err = s.Queries().GetSymbolsByType(c, db.GetSymbolsByTypeParams{
			Type:  symbolType,
			Limit: int32(count),
		})
		if err != nil {
			c.Error(err)
			return
		}
		totalCount, err = s.Queries().GetSymbolCountByType(c, symbolType)
		if err != nil {
			c.Error(err)
			return
		}
	}

	symbolsResp := []symbolResponse{}

	for _, symbol := range symbols {
		symbolsResp = append(symbolsResp, makeSymbolResponse(symbol))
	}

	resp := paginatedSymbols{
		Symbols:    symbolsResp,
		TotalCount: totalCount,
		Pages:      int(totalCount / count),
		ActivePage: 600,
	}

	c.JSON(http.StatusOK, resp)
}

func (s *SymbolHandler) searchSymbol(c *gin.Context) {
	search := c.Query("query")
	count, err := strconv.ParseInt(c.Query("count"), 10, 32)
	if err != nil {
		httputil.AbortBadRequest(c, "param 'count' must be int")
		return
	}

	symbols, err := s.Queries().SearchSymbol(c, db.SearchSymbolParams{
		Search: "%" + search + "%",
		Count:  int32(count),
	})
	if err != nil {
		c.Error(err)
		return
	}

	resp := []symbolResponse{}

	for _, symbol := range symbols {
		resp = append(resp, makeSymbolResponse(symbol))
	}

	c.JSON(http.StatusOK, resp)
}
