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

	file, err := ioutil.ReadFile("data/coingecko/coins-list.json")
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

		exists, err := queries.SymbolExists(ctx, symbol.CoingeckoID)
		if err != nil {
			return fmt.Errorf("could not check if symbol %s exists: %w", symbol.CoingeckoID, err)
		}

		if exists {
			err = queries.UpdateSymbol(ctx, db.UpdateSymbolParams{
				SymbolID:  symbol.CoingeckoID,
				Type:      "crypto",
				Source:    "coingecko",
				Precision: 8,
				SymbolName: sql.NullString{
					String: symbol.Name,
					Valid:  true,
				},
			})
			if err != nil {
				return fmt.Errorf("could not update coingecko coin %s: %w", symbol.CoingeckoID, err)
			}
		} else {
			err = queries.AddSymbol(ctx, db.AddSymbolParams{
				SymbolID:  symbol.CoingeckoID,
				Type:      "crypto",
				Source:    "coingecko",
				Precision: 8,
				SymbolName: sql.NullString{
					String: symbol.Name,
					Valid:  true,
				},
			})
			if err != nil {
				return fmt.Errorf("could not save coingecko coin %s: %w", symbol.CoingeckoID, err)
			}
		}
	}
	return tx.Commit()
}
