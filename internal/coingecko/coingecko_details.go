package coingecko

import (
	"div-dash/internal/db"
	"div-dash/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *CoingeckoService) GetDetails(asset db.Symbol) (model.SymbolDetails, error) {
	resp, err := c.client.R().
		SetPathParam("coin", asset.SymbolID).
		Get("/coins/{coin}")

	if err != nil {
		return model.SymbolDetails{}, fmt.Errorf("could not fetch coin details for %s: %w", asset.SymbolID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return model.SymbolDetails{}, fmt.Errorf("got non-ok status code %d when fetching coin details for %s", resp.StatusCode(), asset.SymbolID)
	}

	var result CoingeckoDetails
	err = json.Unmarshal(resp.Body(), &result)

	if err != nil {
		return model.SymbolDetails{}, fmt.Errorf("could not unmarshal coin details for coin %s: %w", asset.SymbolID, err)
	}

	tags := []model.SymbolTag{
		{
			Label: fmt.Sprintf("#%d Market Cap", result.MarketCapRank),
			Type:  "CHIP",
		},
	}
	if result.HashingAlgorithm != "" {
		tags = append(tags,
			model.SymbolTag{
				Label: result.HashingAlgorithm,
				Type:  "CHIP",
			})
	}

	tags = append(tags,
		model.SymbolTag{
			Label: result.Links.Homepage[0],
			Type:  "LINK",
			Link:  result.Links.Homepage[0],
		})

	return model.SymbolDetails{
		Type:        "crypto",
		Name:        result.Name,
		Description: result.Description.En,
		Tags:        tags,
		Dates:       []model.SymbolDate{},
	}, nil

}
