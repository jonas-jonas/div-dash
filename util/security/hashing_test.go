package security

import (
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
	hashedPassword := "$2a$10$7V7XXwThPDseRqjdkZB8eeTlJKimUVZ2H0/nQ97w4/ri0brf3PYau"

	result := VerifyPassword(password, hashedPassword)

	assert.True(t, result)

}
