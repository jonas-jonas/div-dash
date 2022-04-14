package coingecko

import (
	"context"
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

const COINGECKO_EXCHANGE = "coingecko"

func (c *CoingeckoService) ensureCoingeckoExchange(ctx context.Context) error {
	exists, err := c.Queries().DoesExchangeExist(ctx, COINGECKO_EXCHANGE)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	err = c.Queries().CreateExchange(ctx, db.CreateExchangeParams{
		Exchange:       COINGECKO_EXCHANGE,
		Region:         "", // TODO: Change to null
		Description:    "Coingecko",
		ExchangeSuffix: "",
		Mic:            "",
	})
	return err
}

func (c *CoingeckoService) ImportCryptoSymbols(ctx context.Context) error {
	symbols, err := c.getSymbols()
	if err != nil {
		return fmt.Errorf("could not retrieve coingecko symbols: %w", err)
	}

	err = c.ensureCoingeckoExchange(ctx)

	if err != nil {
		return fmt.Errorf("could not ensure that coingecko exchange exists %w", err)
	}

	const BATCH_SIZE = 1000

	var symbolIDs []string
	var types []string
	var sources []string
	var precisions []int32
	var symbolNames []string

	var exchanges []string
	var exchangeSymbols []string

	tx, err := c.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	queries := c.Queries()

	for index, symbol := range symbols {

		symbolIDs = append(symbolIDs, symbol.ID)
		types = append(types, "crypto")
		sources = append(sources, "coingecko")
		precisions = append(precisions, 4)
		symbolNames = append(symbolNames, symbol.Name)

		exchanges = append(exchanges, "coingecko")
		exchangeSymbols = append(exchangeSymbols, symbol.ID)

		if index%BATCH_SIZE == 0 || index == len(symbols)-1 {

			importParams := db.BulkImportSymbolParams{
				SymbolIds:   symbolIDs,
				Types:       types,
				Sources:     sources,
				Precisions:  precisions,
				SymbolNames: symbolNames,
			}

			importExchangeParams := db.BulkImportSymbolExchangeParams{
				SymbolIds: symbolIDs,
				Types:     types,
				Sources:   sources,
				Exchanges: exchanges,
				Symbols:   exchangeSymbols,
			}

			err = c.executeImport(ctx, queries, importParams, importExchangeParams)
			if err != nil {
				return err
			}

			symbolIDs = nil
			types = nil
			sources = nil
			precisions = nil
			symbolNames = nil

			exchanges = nil
			exchangeSymbols = nil
		}
	}

	return tx.Commit()
}

func (c *CoingeckoService) executeImport(ctx context.Context, queries *db.Queries, importParams db.BulkImportSymbolParams, importExchangeParams db.BulkImportSymbolExchangeParams) error {
	err := queries.BulkImportSymbol(ctx, importParams)

	if err != nil {
		return fmt.Errorf("could not bulk import symbol: %w", err)
	}

	err = queries.BulkImportSymbolExchange(ctx, importExchangeParams)

	if err != nil {
		return fmt.Errorf("could not bulk import symbol with exchange: %w", err)
	}

	return nil
}
