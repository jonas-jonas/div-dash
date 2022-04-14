package registry

import (
	"context"
	"div-dash/internal/account"
	"div-dash/internal/balance"
	"div-dash/internal/coingecko"
	"div-dash/internal/config"
	"div-dash/internal/csvimport"
	"div-dash/internal/db"
	"div-dash/internal/error"
	"div-dash/internal/id"
	"div-dash/internal/identity"
	"div-dash/internal/iex"
	"div-dash/internal/logging"
	"div-dash/internal/symbol"
	"div-dash/internal/timex"
	"div-dash/internal/token"
	"div-dash/internal/transaction"

	"github.com/gin-gonic/gin"
)

type Registry interface {
	db.QueriesProvider
	db.DBProvider
	config.ConfigProvider
	logging.LoggerProvider

	error.ErrorHandlerProvider

	identity.HandlerProvider
	transaction.TransactionHandlerProvider
	symbol.SymbolHandlerProvider
	account.AccountHandlerProvider
	balance.BalanceHandlerProvider

	timex.TimeHolderProvider

	token.TokenServiceProvider

	id.IdServiceProvider

	csvimport.CSVImporterProvider

	iex.IEXServiceProvider

	coingecko.CoingeckoServiceProvider

	RegisterPublicRoutes(ctx context.Context, routes gin.IRoutes)
	RegisterProtectedRoutes(ctx context.Context, routes gin.IRoutes)
	RegisterProtectedMiddleware(ctx context.Context, routes gin.IRoutes)
}
