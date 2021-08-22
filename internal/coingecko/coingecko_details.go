package coingecko

import (
	"div-dash/internal/db"
	"div-dash/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func assembleIndicators(details CoingeckoDetails) []model.SymbolIndicator {
	return []model.SymbolIndicator{
		{
			Label:  "Market Cap",
			Format: "$0.00 a",
			Value:  float64(details.MarketData.MarketCap.Usd),
		},
		{
			Label:  "Circulating Supply",
			Format: "0.00 a",
			Value:  details.MarketData.CirculatingSupply,
		},
		{
			Label:  "Max Supply",
			Format: "0.00 a",
			Value:  details.MarketData.MaxSupply,
		},
		{
			Label:  "Total Supply",
			Format: "0.00 a",
			Value:  details.MarketData.TotalSupply,
		},
	}
}

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
			Label: fmt.Sprintf("Rank #%d", result.MarketCapRank),
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
		Indicators:  assembleIndicators(result),
		Images: model.SymbolImages{
			Thumb: result.Image.Thumb,
		},
		Dates: []model.SymbolDate{},
	}, nil

}
