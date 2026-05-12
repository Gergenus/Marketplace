package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrClaimsFailed = errors.New("claims failed")
	ErrTokenExpired = errors.New("token expired")
	ErrConversion   = errors.New("conversion error")
)

func ParseToken(token string) (string, string, error) {
	keyfunc := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrUnauthorized
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	}

	tkn, err := jwt.Parse(token, keyfunc)
	if err != nil {
		return "", "", err
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", ErrClaimsFailed
	}

	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		return "", "", ErrTokenExpired
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", "", ErrConversion
	}
	uuid, ok := claims["uuid"].(string)
	if !ok {
		return "", "", ErrConversion
	}

	return role, uuid, nil
}
