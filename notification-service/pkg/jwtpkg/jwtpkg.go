package jwtpkg

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTpkg struct {
	JWTMailSecret string
	TokenTTL      time.Duration
}

type JWTpkgInterface interface {
	GenerateToken(email string) (string, error)
}

func NewJWTpkg(Secret string, TokenTTL time.Duration) JWTpkg {
	return JWTpkg{JWTMailSecret: Secret, TokenTTL: TokenTTL}
}

func (j JWTpkg) GenerateToken(email string) (string, error) {
	const op = "jwtpkg.GenerateToken"
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(j.TokenTTL).Unix(),
	})

	token, err := tkn.SignedString([]byte(j.JWTMailSecret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
