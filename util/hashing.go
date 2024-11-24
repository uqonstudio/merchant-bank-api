package util

import "golang.org/x/crypto/bcrypt"

func Encrypt(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func ComparePassword(passwordPayload string, passwordDB string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(passwordPayload))
}
