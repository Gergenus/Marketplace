package jwtpkg

import (
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTinterface interface {
	ParseToken(token string) (string, string, error)
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
	ErrConversion   = errors.New("conversion error")
)

// returns role, uuid and error
func (j JWTpkg) ParseToken(token string) (string, string, error) {
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
		j.log.Error("convesion error", slog.String("error", "claims[role].(string)"))
		return "", "", ErrConversion
	}
	uuid, ok := claims["uuid"].(string)
	if !ok {
		j.log.Error("convesion error", slog.String("error", "claims[uuid].(string)"))
		return "", "", ErrConversion
	}

	return role, uuid, nil
}
