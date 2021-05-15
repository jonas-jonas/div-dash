package services

import (
	"div-dash/internal/config"
	"div-dash/internal/id"
	"div-dash/internal/job"
	"div-dash/internal/mail"
	"div-dash/internal/token"
	"sync"

	"github.com/spf13/viper"
)

var services struct {
	TokenService *token.TokenService
	MailService  *mail.MailService
	IdService    *id.IdService
	JobService   *job.JobService
}

var (
	onceTokenService sync.Once
	onceMailService  sync.Once
	onceIdService    sync.Once
	onceJobService   sync.Once
)

func initTokenService() {
	services.TokenService = token.NewPasetoService(config.Config().Paseto)
}

func TokenService() *token.TokenService {
	onceTokenService.Do(initTokenService)

	return services.TokenService
}

func initMailService() {
	password := viper.GetString("smtp.password")
	port := viper.GetInt("smtp.port")
	server := viper.GetString("smtp.server")
	services.MailService = mail.NewMailService(password, server, port)
}

func MailService() *mail.MailService {
	onceMailService.Do(initMailService)
	return services.MailService
}

func initIdService() {
	services.IdService = id.New()
}

func IdService() *id.IdService {
	onceIdService.Do(initIdService)
	return services.IdService
}


func initJobService() {
	services.JobService = job.New(config.Queries(), config.Logger())
}

func JobService() *job.JobService {
	onceJobService.Do(initJobService)
	return services.JobService
}
