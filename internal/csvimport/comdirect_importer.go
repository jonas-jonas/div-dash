package csvimport

import (
	"context"
	"database/sql"
	"div-dash/internal/db"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/extrame/xls"
	"github.com/shopspring/decimal"
)

const (
	ColumnSettlementDate = 0
	ColumnExecutionDate  = 1
	ColumnWKN            = 2
	ColumnISIN           = 3
	ColumnDescription    = 4
	ColumnSide           = 5
	ColumnAmount         = 6
	ColumnPrice          = 7
	ColumnCurrency       = 8
	ColumnPriceEUR       = 9
	ColumnTotalEUR       = 10
	ColumnFees           = 11
	ColumnOrderID        = 26
)

const (
	ValueSideBuy  = "Kauf"
	ValueSideSell = "Verkauf"
)

const (
	DateFormat = "02.01.2006"
)

var hundred = decimal.New(100, 0)

func (c *csvImporter) importComdirectStatement(ctx context.Context, file multipart.File, accountId string, userId string) error {
	fileReader, err := xls.OpenReader(file, "utf-8")
	if err != nil {
		return fmt.Errorf("could not open excelize reader on file: %w", err)
	}

	sheet := fileReader.GetSheet(0)
	if sheet == nil {
		return fmt.Errorf("no sheet found in file")
	}

	tx, err := c.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	queries := c.Queries().WithTx(tx)

	for i := 1; i <= (int(sheet.MaxRow)); i++ {

		row := sheet.Row(i)
		dateString := row.Col(ColumnSettlementDate)
		wkn := row.Col(ColumnWKN)
		sideString := row.Col(ColumnSide)
		amountString := row.Col(ColumnAmount)
		priceString := row.Col(ColumnPriceEUR)
		orderId := row.Col(ColumnOrderID)

		date, err := time.Parse(DateFormat, dateString)
		if err != nil {
			return fmt.Errorf("unexpected date format '%s' in line %d", dateString, i)
		}

		amount, err := decimal.NewFromString(strings.ReplaceAll(amountString, ",", "."))
		if err != nil {
			return fmt.Errorf("unexpected amount format '%s' in line %d", amountString, i)
		}

		price, err := decimal.NewFromString(strings.ReplaceAll(priceString, ",", "."))
		if err != nil {
			return fmt.Errorf("unexpected price format '%s' in line %d", priceString, i)
		}

		var side string
		if sideString == ValueSideBuy {
			side = db.TransactionSideBuy
		} else if sideString == ValueSideSell {
			side = db.TransactionSideSell
		} else {
			return fmt.Errorf("unknown side '%s' for transaction '%s'", sideString, orderId)
		}

		exists, err := queries.TransactionExists(ctx, sql.NullString{
			String: orderId, Valid: true,
		})

		if err != nil {
			return fmt.Errorf("could not check if transaction '%s' exists: %w", orderId, err)
		}

		if exists {
			log.Printf("skipping transaction '%s' because it already exists", orderId)
			continue
		}

		symbol, err := c.Queries().GetSymbolByWKN(ctx, sql.NullString{
			String: wkn,
			Valid:  true,
		})

		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("unknown WKN '%s'", wkn)
			}
			return fmt.Errorf("could not get symbol of wkn '%s': %w", wkn, err)
		}

		params := db.CreateTransactionParams{
			ID:                  "T" + c.IdService().NewID(5),
			Symbol:              symbol.SymbolID,
			Type:                symbol.Type,
			TransactionProvider: "NONE",
			Price:               price.Mul(hundred).BigInt().Int64(),
			Date:                date,
			Amount:              amount,
			Side:                side,
			AccountID:           accountId,
			UserID:              userId,
			ExternalID: sql.NullString{
				String: orderId,
				Valid:  true,
			},
		}

		_, err = queries.CreateTransaction(ctx, params)

		if err != nil {
			return fmt.Errorf("could not create transaction: %w", err)
		}

	}

	return tx.Commit()
}
