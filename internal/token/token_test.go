package token

import (
	"testing"
	"time"

	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/assert"
)

var key = []byte("YELLOW SUBMARINE, BLACK WIZARDRY")

func createToken(userId string, now, exp, nbt time.Time) (string, error) {

	jsonToken := paseto.JSONToken{
		Audience:   "test",
		Issuer:     "test",
		Jti:        "123",          // TODO
		Subject:    "test_subject", // TODO
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}
	jsonToken.Set("userId", userId)
	footer := "some footer" // TODO
	return paseto.NewV2().Encrypt(key, jsonToken, footer)
}

func TestGenerateToken(t *testing.T) {
	tokenService := New("test_audience", "test_issuer", key)

	token, err := tokenService.GenerateToken(0)

	assert.Nil(t, err)
	assert.Regexp(t, `v2\.local\..*`, token)
}

func TestVerifyToken(t *testing.T) {
	tokenService := New("test_audience", "test_issuer", key)

	token, _ := tokenService.GenerateToken(0)

	result, userId, err := tokenService.VerifyToken(token)

	assert.True(t, result)
	assert.Equal(t, 0, userId)
	assert.Nil(t, err)
}

func TestVerifyTokenWithExpiredToken(t *testing.T) {
	tokenService := New("test_audience", "test_issuer", key)
	now := time.Now().Add(-25 * time.Hour)
	exp := now.Add(24 * time.Hour)
	nbt := now

	token, _ := createToken("1", now, exp, nbt)

	result, userId, err := tokenService.VerifyToken(token)

	assert.False(t, result)
	assert.LessOrEqual(t, -1, userId)
	assert.NotNil(t, err)
	assert.Equal(t, "token expired", err.Error())
}

func TestVerifyTokenWithStringAsUserId(t *testing.T) {
	tokenService := New("test_audience", "test_issuer", key)

	now := time.Now()
	exp := now.Add(24 * time.Hour)
	nbt := now
	token, _ := createToken("userIdString", now, exp, nbt)

	result, userId, err := tokenService.VerifyToken(token)

	assert.False(t, result)
	assert.LessOrEqual(t, -1, userId)
	assert.NotNil(t, err)
	assert.Equal(t, `strconv.Atoi: parsing "userIdString": invalid syntax`, err.Error())
}

func TestVerifyTokenWithDifferentKey(t *testing.T) {
	tokenService := New("test_audience", "test_issuer", key)

	token, _ := tokenService.GenerateToken(0)

	tokenService = New("test_audience", "test_issuer", []byte("111111 SUBMARINE, BLACK WIZARDRY"))

	result, userId, err := tokenService.VerifyToken(token)

	assert.False(t, result)
	assert.LessOrEqual(t, -1, userId)
	assert.NotNil(t, err)
	assert.Equal(t, `invalid token authentication`, err.Error())
}
