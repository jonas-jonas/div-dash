package db

import "fmt"

type TransactionProvider string

const (
	TransactionProviderBinance = "binance"
)

func (e *TransactionProvider) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TransactionProvider(s)
	case string:
		*e = TransactionProvider(s)
	default:
		return fmt.Errorf("unsupported scan type for TransactionProvider: %T", src)
	}
	return nil
}

type TransactionSide string

const (
	TransactionSideSell = "sell"
	TransactionSideBuy  = "buy"
)

type TransactionType string

const (
	TransactionTypeCrypto TransactionType = "crypto"
)

func (e *TransactionType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TransactionType(s)
	case string:
		*e = TransactionType(s)
	default:
		return fmt.Errorf("unsupported scan type for TransactionType: %T", src)
	}
	return nil
}

type UserStatus string

const (
	UserStatusRegistered  = "registered"
	UserStatusActivated   = "activated"
	UserStatusDeactivated = "deactivated"
)

func (e *UserStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserStatus(s)
	case string:
		*e = UserStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for UserStatus: %T", src)
	}
	return nil
}
