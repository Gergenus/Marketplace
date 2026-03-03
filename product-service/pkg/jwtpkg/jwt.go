package jwtpkg

import (
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTinterface interface {
	ParseToken(token string) (string, int, error)
}

type JWTpkg struct {
	Secret string
	log    *slog.Logger
}

func NewJWTpkg(Secret string, log *slog.Logger) JWTpkg {
	return JWTpkg{
		Secret: Secret,
		log:    log,
	}
}

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrClaimsFailed = errors.New("claims failed")
	ErrTokenExpired = errors.New("token expired")
)

func (j JWTpkg) ParseToken(token string) (string, int, error) {
	keyfunc := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrUnauthorized
		}
		return []byte(j.Secret), nil
	}

	tkn, err := jwt.Parse(token, keyfunc)
	if err != nil {
		j.log.Error("error parsing token", slog.String("error", err.Error()))
		return "", -1, err
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return "", -1, ErrClaimsFailed
	}

	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		return "", -1, ErrTokenExpired
	}

	role := claims["role"].(string)
	sellerID := claims["seller_id"].(float64)
	return role, int(sellerID), nil
}
