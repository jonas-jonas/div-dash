package coingecko

import (
	"context"
	"div-dash/internal/db"
	"div-dash/internal/model"
	"encoding/json"
	"fmt"
	"time"
)

func (c *CoingeckoService) GetChart(ctx context.Context, asset db.Symbol, span int) (model.Chart, error) {

	symbol, err := c.queries.GetSymbolOfSymbolAndExchange(ctx, db.GetSymbolOfSymbolAndExchangeParams{
		SymbolID: asset.SymbolID,
		Exchange: "coingecko",
	})

	if err != nil {
		return model.Chart{}, fmt.Errorf("could not get exchanges of asset %s: %w", asset.SymbolID, err)
	}

	if !symbol.Valid {
		return model.Chart{}, fmt.Errorf("no symbol found for symbol %s on exchange %s", asset.SymbolID, "coingecko")
	}

	resp, err := c.client.R().
		SetPathParam("coin", symbol.String).
		SetQueryParam("vs_currency", "eur").
		SetQueryParam("days", "365").
		Get("/coins/{coin}/market_chart")

	if err != nil {
		return model.Chart{}, fmt.Errorf("could not fetch chart for coin %s: %w", asset.SymbolID, err)
	}

	if resp.StatusCode() != 200 {
		return model.Chart{}, fmt.Errorf("got non-ok status from chart for coin %s: %d", asset.SymbolID, resp.StatusCode())
	}

	var chartResponse CoingeckoChart

	err = json.Unmarshal(resp.Body(), &chartResponse)

	if err != nil {
		return model.Chart{}, fmt.Errorf("could not unmarshal price response: %w", err)
	}

	result := model.Chart{}

	for _, entry := range chartResponse.Prices {
		timestamp := entry[0]
		value := entry[1]
		result = append(result, model.ChartEntry{
			Date:  time.Unix(int64(timestamp)/1000, 0).Format("2006-01-02"),
			Price: value,
		})
	}

	return result, nil
}
