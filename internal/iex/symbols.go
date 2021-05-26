package iex

import (
	"context"
	"database/sql"
	"div-dash/internal/config"
	"div-dash/internal/db"
	"encoding/json"
	"errors"
	"io/ioutil"
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

const IEX_IMPORT_JOB_NAME = "import-iex-symbols"

func (i *IEXService) SaveSymbols() error {
	// TODO: Get all Exchanges and collect all regions
	ctx := context.Background()
	week := 60 * 60 * 24 * 7
	expired, err := i.jobService.HasLastSuccessfulJobExpired(ctx, IEX_IMPORT_JOB_NAME, time.Duration(week))
	if !expired || err == nil {
		return errors.New("import-iex-symbols: last successful import was less than a week ago")
	}

	symbols, err := i.getSymbolsByRegion("de")
	if err != nil {
		return err
	}

	count := len(symbols)
	config.Logger().Printf("Importing %v IEX Assets...", count)
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := i.queries.WithTx(tx)

	job, err := i.jobService.StartJob(ctx, IEX_IMPORT_JOB_NAME)
	if err != nil {
		return err
	}

	for _, symbol := range symbols {
		symbolId := symbol.Symbol
		if strings.Contains(symbol.Symbol, "-") {
			parts := strings.Split(symbol.Symbol, "-")
			symbolId = parts[0]
		}
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
			config.Logger().Printf("Could not save iex symbol %s: %s", symbol.Symbol, err.Error())
			continue
		}
		err = queries.ConnectSymbolWithExchange(ctx, db.ConnectSymbolWithExchangeParams{
			Symbol:   symbolId,
			Exchange: symbol.Exchange,
		})

		if err != nil {
			config.Logger().Printf("Could not save iex symbol %s: %s", symbol.Symbol, err.Error())
			continue
		}
	}
	err = tx.Commit()
	if err != nil {
		i.jobService.FailJob(ctx, job.ID, err.Error())
		return err
	}
	i.jobService.FinishJob(ctx, job.ID)
	return nil
}
