package token

import (
	"div-dash/internal/config"

	"pkg.re/essentialkaos/branca.v1"
)

type TokenService struct {
	key        []byte
	tokenValid uint32
}

func NewTokenService(config config.TokenConfiguration) *TokenService {
	return &TokenService{[]byte(config.Key), config.TokenValid}
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
