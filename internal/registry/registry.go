package registry

import (
	"context"
	"div-dash/internal/account"
	"div-dash/internal/config"
	"div-dash/internal/db"
	"div-dash/internal/error"
	"div-dash/internal/identity"
	"div-dash/internal/logging"
	"div-dash/internal/symbol"
	"div-dash/internal/timex"
	"div-dash/internal/token"
	"div-dash/internal/transaction"

	"github.com/gin-gonic/gin"
)

type Registry interface {
	db.QueriesProvider
	config.ConfigProvider
	logging.LoggerProvider

	error.ErrorHandlerProvider

	identity.HandlerProvider
	transaction.TransactionHandlerProvider
	symbol.SymbolHandlerProvider
	account.AccountHandlerProvider

	timex.TimeHolderProvider

	token.TokenServiceProvider

	RegisterPublicRoutes(ctx context.Context, routes gin.IRoutes)
	RegisterProtectedRoutes(ctx context.Context, routes gin.IRoutes)
	RegisterProtectedMiddlware(ctx context.Context, routes gin.IRoutes)
}
