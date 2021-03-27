package services

import (
	"div-dash/internal/token"
	"div-dash/internal/user"
	"sync"
)

var services struct {
	TokenService *token.TokenService
	UserService  *user.UserService
}

var onceTokenService sync.Once
var onceUserService sync.Once

func initTokenService() {
	services.TokenService = token.New(
		"audience", "issuer", []byte("YELLOW SUBMARINE, BLACK WIZARDRY"))
}

func TokenService() *token.TokenService {
	onceTokenService.Do(initTokenService)

	return services.TokenService
}

func initUserService() {
	services.UserService = user.New()
}

func UserService() *user.UserService {
	onceUserService.Do(initUserService)

	return services.UserService
}
