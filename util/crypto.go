package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(value string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(password), nil
}

func CheckPassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("password does not match: %w", err)
	}

	return nil
}
