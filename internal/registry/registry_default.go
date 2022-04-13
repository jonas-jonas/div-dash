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
	"go.uber.org/zap"
)

type RegistryDefault struct {
	queries *db.Queries
	config  *config.Configuration
	logger  *zap.SugaredLogger

	errorHandler *error.Handler

	identity           *identity.Handler
	transactionHandler *transaction.TransactionHandler
	symbolHandler      *symbol.SymbolHandler
	accountHandler     *account.Handler
	tokenService       *token.TokenService
}

func NewRegistryDefault() *RegistryDefault {
	return &RegistryDefault{}
}

func (r *RegistryDefault) RegisterPublicRoutes(ctx context.Context, routes gin.IRoutes) {
	r.IdentityHandler().RegisterPublicRoutes(routes)
}

func (r *RegistryDefault) RegisterProtectedRoutes(ctx context.Context, routes gin.IRoutes) {
	r.IdentityHandler().RegisterPrivateRoutes(routes)
	r.TransactionHandler().RegisterProtectedRoutes(routes)
	r.SymbolHandler().RegisterProtectedRoutes(routes)
}

func (r *RegistryDefault) RegisterProtectedMiddleware(ctx context.Context, routes gin.IRoutes) {
	r.IdentityHandler().RegisterMiddleware(routes)
}

func (r *RegistryDefault) Queries() *db.Queries {
	if r.queries == nil {
		r.queries = db.NewQueries(r)
	}
	return r.queries
}

func (r *RegistryDefault) Config() *config.Configuration {
	if r.config == nil {
		r.config = config.NewConfig()
	}
	return r.config
}

func (r *RegistryDefault) Logger() *zap.SugaredLogger {
	if r.logger == nil {
		r.logger = logging.NewLogger()
	}
	return r.logger
}

func (r *RegistryDefault) ErrorHandler() *error.Handler {
	if r.errorHandler == nil {
		r.errorHandler = error.NewHandler(r)
	}
	return r.errorHandler
}

func (r *RegistryDefault) IdentityHandler() *identity.Handler {
	if r.identity == nil {
		r.identity = identity.NewHandler(r)
	}
	return r.identity
}

func (r *RegistryDefault) TransactionHandler() *transaction.TransactionHandler {
	if r.transactionHandler == nil {
		r.transactionHandler = transaction.NewTransactionHandler(r)
	}
	return r.transactionHandler
}

func (r *RegistryDefault) SymbolHandler() *symbol.SymbolHandler {
	if r.symbolHandler == nil {
		r.symbolHandler = symbol.NewSymbolHandler(r)
	}
	return r.symbolHandler
}

func (r *RegistryDefault) AccountHandler() *account.Handler {
	if r.accountHandler == nil {
		r.accountHandler = account.NewAccountHandler(r)
	}
	return r.accountHandler
}

func (r *RegistryDefault) TimeHolder() timex.TimeHolder {
	return timex.NewTimeProvider()
}

func (r *RegistryDefault) TokenService() *token.TokenService {
	if r.tokenService == nil {
		r.tokenService = token.NewTokenService(r)
	}
	return r.tokenService
}
