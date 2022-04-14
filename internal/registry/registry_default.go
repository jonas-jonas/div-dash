package registry

import (
	"context"
	"database/sql"
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
	"div-dash/internal/mail"
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
	balanceHandler     *balance.BalanceHandler

	tokenService *token.TokenService
	mailService  *mail.MailService
	db           *sql.DB
	csvImporter  csvimport.CSVImporter

	iexService       *iex.IEXService
	coingeckoService *coingecko.CoingeckoService
}

func NewRegistryDefault() Registry {
	return &RegistryDefault{}
}

func (r *RegistryDefault) RegisterPublicRoutes(ctx context.Context, routes gin.IRoutes) {
	r.IdentityHandler().RegisterPublicRoutes(routes)
}

func (r *RegistryDefault) RegisterProtectedRoutes(ctx context.Context, routes gin.IRoutes) {
	r.IdentityHandler().RegisterPrivateRoutes(routes)
	r.TransactionHandler().RegisterProtectedRoutes(routes)
	r.SymbolHandler().RegisterProtectedRoutes(routes)
	r.BalanceHandler().RegisterProtectedRoutes(routes)
	r.AccountHandler().RegisterProtectedRoutes(routes)
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

func (r *RegistryDefault) MailService() *mail.MailService {
	if r.mailService == nil {
		r.mailService = mail.NewMailService(r)
	}
	return r.mailService
}

func (r *RegistryDefault) IdService() id.IdService {
	return id.NewIdService()
}

func (r *RegistryDefault) DB() *sql.DB {
	if r.db == nil {
		r.db = db.NewDB(r)
	}
	return r.db
}

func (r *RegistryDefault) CSVImporter() csvimport.CSVImporter {
	if r.csvImporter == nil {
		r.csvImporter = csvimport.NewCSVImporter(r)
	}
	return r.csvImporter
}

func (r *RegistryDefault) IEXService() *iex.IEXService {
	if r.iexService == nil {
		r.iexService = iex.New(r)
	}
	return r.iexService
}

func (r *RegistryDefault) CoingeckoService() *coingecko.CoingeckoService {
	if r.coingeckoService == nil {
		r.coingeckoService = coingecko.New(r)
	}
	return r.coingeckoService
}

func (r *RegistryDefault) BalanceHandler() *balance.BalanceHandler {
	if r.balanceHandler == nil {
		r.balanceHandler = balance.NewBalanceHandler(r)
	}
	return r.balanceHandler
}
