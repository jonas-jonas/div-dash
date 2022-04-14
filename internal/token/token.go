package token

import (
	"div-dash/internal/config"

	"pkg.re/essentialkaos/branca.v1"
)

type (
	tokenServiceDependencies interface {
		config.ConfigProvider
	}

	TokenServiceProvider interface {
		TokenService() *TokenService
	}
	TokenService struct {
		tokenServiceDependencies
		key        []byte
		tokenValid uint32
	}
)

func NewTokenService(t tokenServiceDependencies) *TokenService {
	return &TokenService{
		tokenServiceDependencies: t,

		key:        []byte(t.Config().Token.Key),
		tokenValid: t.Config().Token.TokenValid,
	}
}

func (t *TokenService) GenerateToken(userId string) (string, error) {

	branca, err := branca.NewBranca(t.key)
	if err != nil {
		return "", err
	}

	branca.SetTTL(t.tokenValid)

	return branca.EncodeToString([]byte(userId))
}

func (t *TokenService) VerifyToken(tokenString string) (bool, string, error) {

	branca, err := branca.NewBranca(t.key)
	if err != nil {
		return false, "", err
	}

	token, err := branca.DecodeString(tokenString)
	if err != nil {
		return false, "", err
	}

	idString := string(token.Payload())
	return true, idString, nil

}
