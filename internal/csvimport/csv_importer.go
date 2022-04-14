package csvimport

import (
	"context"
	"div-dash/internal/db"
	"div-dash/internal/id"
	"errors"
	"fmt"
	"mime/multipart"
)

type (
	csvImporterDependencies interface {
		db.QueriesProvider
		id.IdServiceProvider
		db.DBProvider
	}

	CSVImporterProvider interface {
		CSVImporter() CSVImporter
	}

	CSVImporter interface {
		ImportCSV(ctx context.Context, file multipart.File, accountId string, userId string) error
	}

	csvImporter struct {
		csvImporterDependencies
	}
)

func NewCSVImporter(c csvImporterDependencies) CSVImporter {
	return &csvImporter{csvImporterDependencies: c}
}

func (c *csvImporter) ImportCSV(ctx context.Context, file multipart.File, accountId string, userId string) error {

	account, err := c.Queries().GetAccount(ctx, accountId)
	if err != nil {
		return err
	}

	if !account.AccountType.Valid {
		return fmt.Errorf("statement upload not supported for untyped accounts")
	}

	accountType := account.AccountType.String
	if accountType == "comdirect" {
		return c.importComdirectStatement(ctx, file, accountId, userId)
	} else if accountType == "scalable_capital" {
		return c.importScalableCapitalStatement(ctx, file, accountId, userId)
	}

	return errors.New("unknown account type " + accountType + " for account wiht id " + accountId)
}
