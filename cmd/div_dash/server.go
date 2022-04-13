package main

import (
	"context"
	"div-dash/internal/coingecko"
	"div-dash/internal/iex"
	"div-dash/internal/job"
	"div-dash/internal/registry"
	"div-dash/internal/services"
	"time"

	"github.com/gin-contrib/static"
	ginzap "github.com/gin-contrib/zap"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	// log.SetFlags(log.Lshortfile | log.LstdFlags)
	reg := registry.NewRegistryDefault()

	r := gin.New()
	r.Use(ginzap.Ginzap(reg.Logger().Desugar(), time.RFC3339, false))
	r.Use(ginzap.RecoveryWithZap(reg.Logger().Desugar(), true))

	jobService := job.NewJobService(reg)

	jobService.RunJob(iex.IEXExchangesImportJob, services.IEXService().SaveExchanges)
	jobService.RunJob(iex.IEXImportSymbolsJob, services.IEXService().SaveSymbols)
	jobService.RunJob(iex.ISINAndWKNImportJob, services.IEXService().ImportISINAndWKN)
	jobService.RunJob(coingecko.CoingeckoImportCoinsJob, services.CoingeckoService().ImportCryptoSymbols)

	r.Use(static.Serve("/", static.LocalFile("web/build", true)))

	r.NoRoute(func(c *gin.Context) {
		c.File("web/build/index.html")
	})

	api := r.Group("/api")
	api.Use(reg.ErrorHandler().HandleErrors)
	reg.RegisterPublicRoutes(context.TODO(), api)

	authorizedApi := api.Group("/")

	reg.RegisterProtectedMiddleware(context.TODO(), authorizedApi)

	reg.RegisterProtectedRoutes(context.TODO(), api)

	port := viper.GetString("server.port")
	r.Run(":" + port)
}
