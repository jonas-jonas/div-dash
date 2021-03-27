package security

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {

	passwordByte := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
