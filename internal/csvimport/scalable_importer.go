package csvimport

import (
	"context"
	"errors"
	"mime/multipart"
)

func (c *csvImporter) importScalableCapitalStatement(ctx context.Context, file multipart.File, accountId string, userId string) error {
	return errors.New("not implemented")
}
