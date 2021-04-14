package token

import (
	"div-dash/internal/config"
	"testing"
	"time"

	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/assert"
)

var key = "YELLOW SUBMARINE, BLACK WIZARDRY"

func createToken(userId string, now, exp, nbt time.Time) (string, error) {

	jsonToken := paseto.JSONToken{
		Audience:   "test",
		Issuer:     "test",
		Jti:        "123", // TODO
		Subject:    userId,
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}
	footer := "some footer" // TODO
	return paseto.NewV2().Encrypt([]byte(key), jsonToken, footer)
}

func TestGenerateToken(t *testing.T) {
	tokenService := NewPasetoService(config.PasetoConfiguration{
		Audience: "test_audience", Issuer: "test_issuer", Key: key, TokenValid: 24,
	})

	token, err := tokenService.GenerateToken(0)

	assert.Nil(t, err)
	assert.Regexp(t, `v2\.local\..*`, token)
}

func TestVerifyToken(t *testing.T) {
	tokenService := NewPasetoService(config.PasetoConfiguration{
		Audience: "test_audience", Issuer: "test_issuer", Key: key, TokenValid: 24,
	})

	token, _ := tokenService.GenerateToken(0)

	result, userId, err := tokenService.VerifyToken(token)

	assert.True(t, result)
	assert.Equal(t, int64(0), userId)
	assert.Nil(t, err)
}

func TestVerifyTokenWithExpiredToken(t *testing.T) {
	tokenService := NewPasetoService(config.PasetoConfiguration{
		Audience: "test_audience", Issuer: "test_issuer", Key: key, TokenValid: 24,
	})
	now := time.Now().Add(-25 * time.Hour)
	exp := now.Add(24 * time.Hour)
	nbt := now

	token, _ := createToken("1", now, exp, nbt)

	result, userId, err := tokenService.VerifyToken(token)

	assert.False(t, result)
	assert.LessOrEqual(t, int64(-1), userId)
	assert.NotNil(t, err)
	assert.Equal(t, "token expired", err.Error())
}

func TestVerifyTokenWithStringAsUserId(t *testing.T) {
	tokenService := NewPasetoService(config.PasetoConfiguration{
		Audience: "test_audience", Issuer: "test_issuer", Key: key, TokenValid: 24,
	})

	now := time.Now()
	exp := now.Add(24 * time.Hour)
	nbt := now
	token, _ := createToken("userIdString", now, exp, nbt)

	result, userId, err := tokenService.VerifyToken(token)

	assert.False(t, result)
	assert.LessOrEqual(t, int64(-1), userId)
	assert.NotNil(t, err)
	assert.Equal(t, `strconv.ParseInt: parsing "userIdString": invalid syntax`, err.Error())
}

func TestVerifyTokenWithDifferentKey(t *testing.T) {
	tokenService := NewPasetoService(config.PasetoConfiguration{
		Audience: "test_audience", Issuer: "test_issuer", Key: key, TokenValid: 24,
	})

	token, _ := tokenService.GenerateToken(0)

	tokenService = NewPasetoService(config.PasetoConfiguration{
		Audience: "test_audience", Issuer: "test_issuer", Key: "111111 SUBMARINE, BLACK WIZARDRY", TokenValid: 24,
	})

	result, userId, err := tokenService.VerifyToken(token)

	assert.False(t, result)
	assert.LessOrEqual(t, int64(-1), userId)
	assert.NotNil(t, err)
	assert.Equal(t, `invalid token authentication`, err.Error())
}
