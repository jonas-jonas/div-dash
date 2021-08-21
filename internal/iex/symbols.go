package iex

import (
	"context"
	"database/sql"
	"div-dash/internal/db"
	"div-dash/internal/job"
	"encoding/json"
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
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := i.queries.WithTx(tx)

	for _, symbol := range symbols {
		symbolId := symbol.Symbol
		if strings.Contains(symbol.Symbol, "-") {
			parts := strings.Split(symbol.Symbol, "-")
			symbolId = parts[0]
		}

		exists, err := queries.SymbolExists(ctx, symbolId)
		if err != nil {
			log.Printf("Could not check if symbol %s exists: %s", symbolId, err.Error())
			return err
		}

		if exists {
			err = queries.UpdateSymbol(ctx, db.UpdateSymbolParams{
				SymbolID: symbolId,
				SymbolName: sql.NullString{
					Valid:  true,
					String: symbol.Name,
				},
				Source:    "iex",
				Type:      symbol.Type,
				Precision: 4,
			})
			if err != nil {
				log.Printf("Could not update iex symbol %s: %s", symbol.Symbol, err.Error())
				continue
			}

		} else {
			err = queries.AddSymbol(ctx, db.AddSymbolParams{
				SymbolID: symbolId,
				SymbolName: sql.NullString{
					Valid:  true,
					String: symbol.Name,
				},
				Source:    "iex",
				Type:      symbol.Type,
				Precision: 4,
			})

			if err != nil {
				log.Printf("Could not save iex symbol %s: %s", symbol.Symbol, err.Error())
				continue
			}
		}
		err = queries.ConnectSymbolWithExchange(ctx, db.ConnectSymbolWithExchangeParams{
			Symbol:   symbolId,
			Exchange: symbol.Exchange,
		})

		if err != nil {
			log.Printf("Could not save iex symbol %s: %s", symbol.Symbol, err.Error())
			continue
		}
	}
	return tx.Commit()
}
