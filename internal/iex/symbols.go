package iex

import (
	"context"
	"div-dash/internal/db"
	"div-dash/internal/job"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Symbol struct {
	Symbol         string `json:"symbol"`
	Exchange       string `json:"exchange"`
	ExchangeSuffix string `json:"exchangeSuffix"`
	ExchangeName   string `json:"exchangeName"`
	Name           string `json:"name"`
	Date           string `json:"date"`
	Type           string `json:"type"`
	Region         string `json:"region"`
	Currency       string `json:"currency"`
	IsEnabled      bool   `json:"isEnabled"`
	Figi           string `json:"figi"`
	Cik            string `json:"cik"`
	Lei            string `json:"lei"`
}

func (i *IEXService) getSymbolsByRegion(region string) ([]Symbol, error) {

	file, err := ioutil.ReadFile("data/iex/Symbols-de.json")
	if err != nil {
		return nil, err
	}

	symbols := []Symbol{}

	err = json.Unmarshal([]byte(file), &symbols)
	if err != nil {
		return nil, err
	}
	return symbols, nil

	// token := "pk_f63a9516a1d14334bcf987d1dd52af64"
	// resp, err := i.client.R().
	// 	SetQueryParam("token", token).
	// 	SetPathParam("region", region).
	// 	Get("https://cloud.iexapis.com/stable/ref-data/region/{region}/symbols")

	// if err != nil {
	// 	return nil, err
	// }

	// if resp.StatusCode() != 200 {
	// 	return nil, fmt.Errorf("iexService/getSymbolsByRegion: error fetching symbols for region %s: %s", region, resp.RawResponse.Proto)
	// }

	// body := resp.Body()
	// symbols := []Symbol{}

	// err = json.Unmarshal(body, &symbols)
	// if err != nil {
	// 	return nil, err
	// }
	// return symbols, nil
}

var IEXImportSymbolsJob job.JobDefinition = job.JobDefinition{
	Key:      "import-iex-symbols",
	Validity: 24 * 7 * time.Hour,
}

func (i *IEXService) SaveSymbols(ctx context.Context) error {

	symbols, err := i.getSymbolsByRegion("de")
	if err != nil {
		return err
	}

	count := len(symbols)
	log.Printf("Importing %v IEX Assets...", count)
	tx, err := i.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := i.Queries().WithTx(tx)

	const BATCH_SIZE = 1000

	var symbolIDs []string
	var types []string
	var sources []string
	var precisions []int32
	var symbolNames []string

	var exchanges []string
	var exchangeSymbols []string

	for index, symbol := range symbols {
		symbolId := symbol.Symbol
		if strings.Contains(symbol.Symbol, "-") {
			parts := strings.Split(symbol.Symbol, "-")
			symbolId = parts[0]
		}

		symbolIDs = append(symbolIDs, symbolId)
		types = append(types, symbol.Type)
		sources = append(sources, "iex")
		precisions = append(precisions, 4)
		symbolNames = append(symbolNames, symbol.Name)

		exchanges = append(exchanges, symbol.Exchange)
		exchangeSymbols = append(exchangeSymbols, symbolId+"-"+symbol.ExchangeSuffix)

		if index%BATCH_SIZE == 0 || index == len(symbols)-1 {
			err = queries.BulkImportSymbol(ctx, db.BulkImportSymbolParams{
				SymbolIds:   symbolIDs,
				Types:       types,
				Sources:     sources,
				Precisions:  precisions,
				SymbolNames: symbolNames,
			})

			if err != nil {
				return fmt.Errorf("could not bulk import symbol %s: %w", symbol.Symbol, err)
			}

			err = queries.BulkImportSymbolExchange(ctx, db.BulkImportSymbolExchangeParams{
				SymbolIds: symbolIDs,
				Types:     types,
				Sources:   sources,
				Exchanges: exchanges,
				Symbols:   exchangeSymbols,
			})

			if err != nil {
				// TODO: Fix error messages here?
				return fmt.Errorf("could not bulk import symbol with exchange %s @ %s: %w", symbol.Symbol, symbol.Exchange, err)
			}

			symbolIDs = nil
			types = nil
			sources = nil
			precisions = nil
			symbolNames = nil

			exchanges = nil
			exchangeSymbols = nil
		}

		if err != nil {
			return fmt.Errorf("could not save iex symbol %s: %s", symbol.Symbol, err.Error())
		}
	}
	return tx.Commit()
}
