package coingecko

import (
	"context"
	"database/sql"
	"div-dash/internal/db"
	"div-dash/internal/job"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

func (c *CoingeckoService) getSymbols() ([]CoingeckoCoin, error) {

	file, err := ioutil.ReadFile("data/coingecko/coins-markets.json")
	if err != nil {
		return nil, err
	}

	symbols := []CoingeckoCoin{}

	err = json.Unmarshal([]byte(file), &symbols)
	if err != nil {
		return nil, err
	}
	return symbols, nil
}

var CoingeckoImportCoinsJob job.JobDefinition = job.JobDefinition{
	Key:      "import-coingecko-coins",
	Validity: 24 * 7 * time.Hour,
}

func (c *CoingeckoService) ImportCryptoSymbols(ctx context.Context) error {
	symbols, err := c.getSymbols()
	if err != nil {
		return fmt.Errorf("could not retrieve coingecko symbols: %w", err)
	}

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	queries := c.queries.WithTx(tx)

	for _, symbol := range symbols {

		exists, err := queries.SymbolExists(ctx, symbol.Symbol)
		if err != nil {
			return fmt.Errorf("could not check if symbol %s exists: %w", symbol.Symbol, err)
		}

		if exists {
			err = queries.UpdateSymbol(ctx, db.UpdateSymbolParams{
				SymbolID:  symbol.Symbol,
				Type:      "crypto",
				Source:    "coingecko",
				Precision: 8,
				SymbolName: sql.NullString{
					String: symbol.Name,
					Valid:  true,
				},
			})
			if err != nil {
				return fmt.Errorf("could not update coingecko coin %s: %w", symbol.Symbol, err)
			}
		} else {
			err = queries.AddSymbol(ctx, db.AddSymbolParams{
				SymbolID:  symbol.Symbol,
				Type:      "crypto",
				Source:    "coingecko",
				Precision: 8,
				SymbolName: sql.NullString{
					String: symbol.Name,
					Valid:  true,
				},
			})
			if err != nil {
				return fmt.Errorf("could not save coingecko coin %s: %w", symbol.ID, err)
			}
		}

		queries.ConnectSymbolWithExchange(ctx, db.ConnectSymbolWithExchangeParams{
			SymbolID: symbol.Symbol,
			Exchange: "coingecko",
			Symbol: sql.NullString{
				String: symbol.ID,
				Valid:  true,
			},
		})
	}
	return tx.Commit()
}
