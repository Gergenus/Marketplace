package service

import (
	"context"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/google/uuid"
)

type UserServiceInterface interface {
	AddUser(ctx context.Context, user models.User) (*uuid.UUID, error)
	Login(ctx context.Context, email, password, userAgent, ip string) (string, string, error)
	RefreshToken(ctx context.Context, oldUuid uuid.UUID, userAgent, ip string, oldAccessToken string) (*uuid.UUID, string, error)
	Logout(ctx context.Context, refreshToken string) error
	Verification(ctx context.Context, token string) error
	ConfirmUser(ctx context.Context, uuid string) error
}
