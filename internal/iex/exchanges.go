package iex

import (
	"context"
	"div-dash/internal/config"
	"div-dash/internal/db"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Exchanges struct {
	Exchange       string `json:"exchange"`
	Region         string `json:"region"`
	Description    string `json:"description"`
	Mic            string `json:"mic"`
	ExchangeSuffix string `json:"exchangeSuffix"`
}

func (i *IEXService) getExchanges() ([]Exchanges, error) {
	// Add Token from config here
	token := "pk_f63a9516a1d14334bcf987d1dd52af64"
	resp, err := i.client.R().
		SetQueryParam("token", token).
		Get("https://cloud.iexapis.com/stable/ref-data/exchanges")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("iexService/getExchanges: error fetching exchanges: %s", resp.Body())
	}

	body := resp.Body()

	exchanges := []Exchanges{}

	err = json.Unmarshal(body, &exchanges)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal exchanges '%s': %w", body, err)
	}
	return exchanges, nil
}

const IEX_IMPORT_EXCHANGES_JOB_NAME = "iex-import-exchanges"

func (i *IEXService) SaveExchanges() error {

	ctx := context.Background()
	week := 60 * 60 * 24 * 7
	expired, err := i.jobService.HasLastSuccessfulJobExpired(ctx, IEX_IMPORT_EXCHANGES_JOB_NAME, time.Duration(week))

	if err != nil {
		return fmt.Errorf("import-iex-exchanges: error while checking last exchange import: %w", err)
	}

	if !expired {
		return errors.New("import-iex-exchanges: last successful exchange import was less than a week ago")
	}

	job, err := i.jobService.StartJob(ctx, IEX_IMPORT_EXCHANGES_JOB_NAME)
	if err != nil {
		return err
	}

	exchanges, err := i.getExchanges()

	if err != nil {
		return err
	}

	count := len(exchanges)
	config.Logger().Printf("Importing %v IEX Exchanges...", count)
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := i.queries.WithTx(tx)

	for _, exchange := range exchanges {

		err = queries.CreateExchange(ctx, db.CreateExchangeParams(exchange))
		if err != nil {
			config.Logger().Printf("Could not save iex exchange %s: %s", exchange.Exchange, err.Error())
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
