package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword is a password hash function using bcrypt
func HashPassword(password string) (string, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("%s", err)
	}

	return string(hashPassword), nil
}

func CompareHashPassword(password, hp string) error {

	return bcrypt.CompareHashAndPassword([]byte(hp), []byte(password))
}
