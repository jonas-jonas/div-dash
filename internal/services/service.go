package services

import (
	"div-dash/internal/mail"
	"div-dash/internal/token"
	"sync"
)

var services struct {
	TokenService *token.TokenService
	MailService  *mail.MailService
}

var onceTokenService sync.Once
var onceMailService sync.Once

func initTokenService() {
	services.TokenService = token.New(
		"audience", "issuer", []byte("YELLOW SUBMARINE, BLACK WIZARDRY"))
}

func TokenService() *token.TokenService {
	onceTokenService.Do(initTokenService)

	return services.TokenService
}

func initMailService() {
	services.MailService = mail.NewMailService("smtpPass", "localhost", 1025)
}

func MailService() *mail.MailService {
	onceMailService.Do(initMailService)
	return services.MailService
}
