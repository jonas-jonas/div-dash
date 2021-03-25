package token

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/o1egl/paseto"
)

type TokenService struct {
	audience string
	issuer   string
	key      []byte
	paseto   paseto.V2
}

func New(audience, issuer string, key []byte) *TokenService {
	return &TokenService{audience: audience, issuer: issuer, key: key, paseto: paseto.V2{}}
}

func (t *TokenService) GenerateToken(userId int) (string, error) {

	userIdString := strconv.Itoa(userId)

	now := time.Now()
	exp := now.Add(24 * time.Hour)
	nbt := now

	jsonToken := paseto.JSONToken{
		Audience:   t.audience,
		Issuer:     t.issuer,
		Jti:        "123",          // TODO
		Subject:    "test_subject", // TODO
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}
	jsonToken.Set("userId", userIdString)
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

	userIdString := newJsonToken.Get("userId")
	fmt.Printf("exp %v", newJsonToken.Expiration)

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		return false, -1, err
	}

	return true, userId, nil
}
