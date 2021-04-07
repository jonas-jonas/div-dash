package token

import (
	"div-dash/internal/config"
	"errors"
	"strconv"
	"time"

	"github.com/o1egl/paseto"
)

type TokenService struct {
	audience   string
	issuer     string
	key        []byte
	tokenValid int
	paseto     paseto.V2
}

func NewPasetoService(config config.PasetoConfiguration) *TokenService {
	paseto := paseto.V2{}
	return &TokenService{config.Audience, config.Issuer, []byte(config.Key), config.TokenValid, paseto}
}

func (t *TokenService) GenerateToken(userId int64) (string, error) {

	userIdString := strconv.FormatInt(userId, 10)

	now := time.Now()
	exp := now.Add(time.Duration(t.tokenValid) * time.Hour)
	nbt := now

	jsonToken := paseto.JSONToken{
		Audience:   t.audience,
		Issuer:     t.issuer,
		Jti:        "123", // TODO
		Subject:    userIdString,
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}
	footer := "some footer" // TODO

	return t.paseto.Encrypt(t.key, jsonToken, footer)
	// token = "v2.local.E42A2iMY9SaZVzt-WkCi45_aebky4vbSUJsfG45OcanamwXwieieMjSjUkgsyZzlbYt82miN1xD-X0zEIhLK_RhWUPLZc9nC0shmkkkHS5Exj2zTpdNWhrC5KJRyUrI0cupc5qrctuREFLAvdCgwZBjh1QSgBX74V631fzl1IErGBgnt2LV1aij5W3hw9cXv4gtm_jSwsfee9HZcCE0sgUgAvklJCDO__8v_fTY7i_Regp5ZPa7h0X0m3yf0n4OXY9PRplunUpD9uEsXJ_MTF5gSFR3qE29eCHbJtRt0FFl81x-GCsQ9H9701TzEjGehCC6Bhw.c29tZSBmb290ZXI"
}

func (t *TokenService) VerifyToken(token string) (bool, int, error) {

	var newJsonToken paseto.JSONToken
	var newFooter string
	err := t.paseto.Decrypt(token, t.key, &newJsonToken, &newFooter)

	if err != nil {
		return false, -1, err
	}

	if newJsonToken.Expiration.Before(time.Now()) {
		return false, -1, errors.New("token expired")
	}

	userIdString := newJsonToken.Subject

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		return false, -1, err
	}

	return true, userId, nil
}
