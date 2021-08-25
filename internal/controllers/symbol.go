package controllers

import (
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type paginatedSymbols struct {
	Symbols    []symbolResponse `json:"symbols"`
	TotalCount int64            `json:"totalCount"`
	Pages      int              `json:"pages"`
	ActivePage int              `json:"activePage"`
}

func GetSymbols(c *gin.Context) {
	symbolType := c.Query("type")

	count, err := strconv.ParseInt(c.Query("count"), 10, 32)
	if err != nil {
		AbortBadRequest(c, "param 'count' must be int")
		return
	}

	var symbols []db.Symbol
	var totalCount int64

	if symbolType == "" {
		symbols, err = config.Queries().GetSymbols(c, int32(count))
		if err != nil {
			c.Error(err)
			return
		}
		totalCount, err = config.Queries().GetSymbolCount(c)
		if err != nil {
			c.Error(err)
			return
		}
	} else {
		symbols, err = config.Queries().GetSymbolsByType(c, db.GetSymbolsByTypeParams{
			Type:  symbolType,
			Limit: int32(count),
		})
		if err != nil {
			c.Error(err)
			return
		}
		totalCount, err = config.Queries().GetSymbolCountByType(c, symbolType)
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
		resp = append(resp, makeSymbolResponse(symbol))
	}

	c.JSON(http.StatusOK, resp)
}

func GetSymbolDetails(c *gin.Context) {
	symbolId := c.Param("symbolId")

	symbol, err := config.Queries().GetSymbol(c, symbolId)
	if err != nil {
		c.Error(err)
		return
	}

	details, err := services.PriceService().GetDetails(c, symbol)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, details)
}

func GetSymbolChart(c *gin.Context) {
	symbolId := c.Param("symbolId")

	symbol, err := config.Queries().GetSymbol(c, symbolId)
	if err != nil {
		c.Error(err)
		return
	}

	chart, err := services.PriceService().GetChart(c, symbol, 1)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, chart)
}
