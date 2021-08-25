package coingecko

import (
	"context"
	"div-dash/internal/db"
	"encoding/json"
	"fmt"
)

func (c *CoingeckoService) GetPrice(ctx context.Context, asset db.Symbol) (float64, error) {

	symbol, err := c.queries.GetSymbolOfSymbolAndExchange(ctx, db.GetSymbolOfSymbolAndExchangeParams{
		SymbolID: asset.SymbolID,
		Exchange: "coingecko",
	})

	if err != nil {
		return -1, fmt.Errorf("could not get exchanges of asset %s: %w", asset.SymbolID, err)
	}

	if !symbol.Valid {
		return -1, fmt.Errorf("no symbol found for symbol %s on exchange %s", asset.SymbolID, "coingecko")
	}

	resp, err := c.client.R().
		SetQueryParam("ids", symbol.String).
		SetQueryParam("vs_currencies", "EUR").
		Get("/simple/price")

	if err != nil {
		return -1, fmt.Errorf("could not fetch simple price for coin %s: %w", asset.SymbolID, err)
	}

	if resp.StatusCode() != 200 {
		return -1, fmt.Errorf("got non-ok status from simple price for coin %s: %d", asset.SymbolID, resp.StatusCode())
	}

	var priceResponse CoingeckoPriceResponse

	err = json.Unmarshal(resp.Body(), &priceResponse)

	if err != nil {
		return -1, fmt.Errorf("could not unmarshal price response: %w", err)
	}

	return priceResponse[symbol.String]["eur"], nil
}
