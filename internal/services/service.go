package services

import (
	"div-dash/internal/token"
	"sync"
)

var services struct {
	TokenService *token.TokenService
}

var onceTokenService sync.Once

func initTokenService() {
	services.TokenService = token.New(
		"audience", "issuer", []byte("YELLOW SUBMARINE, BLACK WIZARDRY"))
}

func TokenService() *token.TokenService {
	onceTokenService.Do(initTokenService)

	return services.TokenService
}
