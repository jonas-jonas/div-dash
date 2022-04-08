// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	UserID      string         `json:"userID"`
	AccountType sql.NullString `json:"accountType"`
}

type AccountType struct {
	AccountType string `json:"accountType"`
	Label       string `json:"label"`
}

type Exchange struct {
	Exchange       string `json:"exchange"`
	ExchangeSuffix string `json:"exchangeSuffix"`
	Region         string `json:"region"`
	Description    string `json:"description"`
	Mic            string `json:"mic"`
}

type JobHistory struct {
	ID           int32          `json:"id"`
	Name         string         `json:"name"`
	Started      int64          `json:"started"`
	Finished     sql.NullInt64  `json:"finished"`
	ErrorMessage sql.NullString `json:"errorMessage"`
}

type Symbol struct {
	SymbolID   string         `json:"symbolID"`
	Type       string         `json:"type"`
	Source     string         `json:"source"`
	Precision  int32          `json:"precision"`
	SymbolName sql.NullString `json:"symbolName"`
	Isin       sql.NullString `json:"isin"`
	Wkn        sql.NullString `json:"wkn"`
}

type SymbolExchange struct {
	SymbolID string         `json:"symbolID"`
	Type     string         `json:"type"`
	Source   string         `json:"source"`
	Exchange string         `json:"exchange"`
	Symbol   sql.NullString `json:"symbol"`
}

type Transaction struct {
	ID                  string          `json:"id"`
	Symbol              string          `json:"symbol"`
	Type                string          `json:"type"`
	TransactionProvider string          `json:"transactionProvider"`
	Price               int64           `json:"price"`
	Date                time.Time       `json:"date"`
	Amount              decimal.Decimal `json:"amount"`
	AccountID           string          `json:"accountID"`
	UserID              string          `json:"userID"`
	Side                string          `json:"side"`
	ExternalID          sql.NullString  `json:"externalID"`
}

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	Status       string `json:"status"`
}

type UserRegistration struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"userID"`
	Timestamp time.Time `json:"timestamp"`
}
