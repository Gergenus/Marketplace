package hash

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	const op = "hash.HashPassword"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return string(hash), nil
}

func CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func HashToken(token string) string {
	const op = "hash.HashPassword"
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

func CheckToken(hashedToken, token string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token)) == nil
}
