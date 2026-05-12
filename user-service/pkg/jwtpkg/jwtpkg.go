package jwtpkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrClaimsFailed         = errors.New("claims failed")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrTokenExpired         = errors.New("token expired")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrConversion           = errors.New("conversion error")
)

type UserJWTpkg struct {
	Secret        string
	AccessTTL     time.Duration
	EmailTTL      time.Duration
	JWTMailSecret string
}

type UserJWTInterface interface {
	GenerateAccessToken(user models.User) (string, error)
	RegenerateToken(oldToken string) (string, error)
	// returns email
	ParseMailToken(token string) (string, error)
	// returns role, uuid and error
	ParseToken(token string) (string, string, error)
}

func NewUserJWTpkg(Secret string, AccessTTL, EmailTTL time.Duration, JWTMailSecret string) UserJWTpkg {
	return UserJWTpkg{Secret: Secret, AccessTTL: AccessTTL, EmailTTL: EmailTTL, JWTMailSecret: JWTMailSecret}
}

func (u UserJWTpkg) GenerateAccessToken(user models.User) (string, error) {
	const op = "jwtpkg.GenerateAccessToken"
	claims := jwt.MapClaims{}
	claims["exp"] = time.Now().Add(u.AccessTTL).Unix()
	claims["uuid"] = user.ID.String()
	claims["verified"] = user.Verified
	claims["role"] = user.Role
	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	AccessTokenString, err := AccessToken.SignedString([]byte(u.Secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return AccessTokenString, nil
}

func (u UserJWTpkg) RegenerateToken(oldToken string) (string, error) {
	const op = "jwtpkg.RegenerateToken"
	token, err := jwt.Parse(oldToken, func(t *jwt.Token) (interface{}, error) { return []byte(u.Secret), nil }, jwt.WithoutClaimsValidation())
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrClaimsFailed
	}
	UUIDString := claims["uuid"].(string)
	parsedUUID, err := uuid.Parse(UUIDString)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	user := models.User{
		ID:       parsedUUID,
		Role:     claims["role"].(string),
		Verified: claims["verified"].(bool),
	}
	return u.GenerateAccessToken(user)
}

// returns email
func (u UserJWTpkg) ParseMailToken(token string) (string, error) {
	const op = "jwtpkg.ParseMailToken"
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(u.JWTMailSecret), nil
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrClaimsFailed
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", ErrClaimsFailed
	}
	if int64(exp) < time.Now().Unix() {
		return "", ErrTokenExpired
	}
	fmt.Println(exp)
	email := claims["email"]
	fmt.Println(email)
	emailString, ok := email.(string)
	if !ok {
		return "", ErrClaimsFailed
	}
	return emailString, nil
}

// returns role, uuid and error
func (j UserJWTpkg) ParseToken(token string) (string, string, error) {
	keyfunc := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrUnauthorized
		}
		return []byte(j.Secret), nil
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
