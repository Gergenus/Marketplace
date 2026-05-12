package repository

import (
	"context"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/google/uuid"
)

type RepositoryInterface interface {
	AddUser(ctx context.Context, user models.User) (*uuid.UUID, error)
	GetUser(ctx context.Context, email string) (*models.User, error)
	GetUserByUUID(ctx context.Context, uuid string) (*models.User, error)
	CreateJWTSession(ctx context.Context, userId string, refreshToken, fingerprint, ip string, expiresIn int64) error
	GetRefreshSession(ctx context.Context, oldUuid string) (*models.RefreshSession, error)
	DeleteSession(ctx context.Context, refreshToken string) error
	Verification(ctx context.Context, email string) error
}
