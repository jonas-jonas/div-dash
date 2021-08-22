package coingecko

import (
	"div-dash/internal/db"
	"encoding/json"
	"fmt"
)

func (c *CoingeckoService) GetPrice(asset db.Symbol) (float64, error) {

	resp, err := c.client.R().
		SetQueryParam("ids", asset.SymbolID).
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

	return priceResponse[asset.SymbolID]["eur"], nil
}
