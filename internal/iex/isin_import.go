package iex

import (
	"context"
	"database/sql"
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/internal/job"
	"encoding/csv"
	"os"
	"time"
)

var ISINAndWKNImportJob job.JobDefinition = job.JobDefinition{
	Key:      "import-isin-and-wkn",
	Validity: 60 * 31 * time.Hour,
}

func (i *IEXService) ImportISINAndWKN(ctx context.Context) error {

	f, err := os.Open("data/xetra/t7-xetr-allTradableInstruments.csv")
	if err != nil {
		return err
	}
	defer f.Close() // this needs to be after the err check

	reader := csv.NewReader(f)
	reader.Comma = ';'

	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for i, line := range lines {
		if i < 3 {
			continue
		}
		if len(line) < 8 {
			continue
		}
		isin := line[3]
		wknPadded := line[6]
		symbolId := line[7]

		wknLength := len(wknPadded)
		wkn := wknPadded[wknLength-6 : wknLength]

		err = config.Queries().AddISINAndWKN(ctx, db.AddISINAndWKNParams{
			Isin:     sql.NullString{String: isin, Valid: true},
			Wkn:      sql.NullString{String: wkn, Valid: true},
			SymbolID: symbolId,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
