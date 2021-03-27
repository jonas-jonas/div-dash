package security

import (
	"div-dash/util/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashing(t *testing.T) {

	password := "pass"

	hash, err := HashPassword(password)

	assert.Nil(t, err)

	assert.Regexp(t, `^\$2[ayb]\$.{56}$`, hash)
}

func TestVerify(t *testing.T) {

	password := "pass"
	hashedPassword := testutil.PasswordHash
	result := VerifyHash(password, hashedPassword)

	assert.True(t, result)

}
