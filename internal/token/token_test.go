package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var key = []byte("YELLOW SUBMARINE, BLACK WIZARDRY")

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

	token := "v2.local.K5I8B91bBmF_2VBUVYTM-rxCunf5nJHB1Rc61n1LvuqiccstQYKZT7NB_TdbAT_28yGP7dTC4LxL8zCiv0y0-qu0rIPfxTGMV2rjj6MOy8EJWuTgbpSwgxRVSI-c2sCkJpMNUI-OfbZ61gcxpe_-HzL6rkLwunsdQoF5uj6XZ3d4ny0zC9oTTZ7EoJuSazn-nJFZviaIyDH3Z0b5RmPbzU0IYrpLmqVMTR1ZNB4fHAxFO60JJRip-Q3slqoL2EGA7AhfIq9H6IRYjAJDKBSdOufmdMICYUSZgIsDxkw.c29tZSBmb290ZXI"

	result, userId, err := tokenService.VerifyToken(token)

	assert.False(t, result)
	assert.LessOrEqual(t, -1, userId)
	assert.NotNil(t, err)
	assert.Equal(t, "token expired", err.Error())
}

func TestVerifyTokenWithStringAsUserId(t *testing.T) {
	tokenService := New("test_audience", "test_issuer", key)

	token := "v2.local.YlQGie-ce9Tv0XIJER2rDmPR8kezPOaNpcgrIoe4C_Vwa3RrDLhm3IlrelDbkDsJ3HAEbHAO6Tby-AvTfY1iItKEelZchLSmV5A-TiePC8HYiCjS1VtVwe4p-PGI2o0JXo10mzGdnY7jtPXW1UUgBT5e4oGWIefQr3sf-xZcDbJlUbIzyL1_ERoVsiPA58DMLiKUqkAgh-0BgW22hCTVr7iJ8wIs1C018SmTq9My1PerTCaL4TnZucpxvtlHLVK6CByjpj786mGXPptTFPmiOaQfpNCvLajFIVlyVzcBJtnBHpTKCfdQzw.c29tZSBmb290ZXI"

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
