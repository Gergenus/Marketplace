package repository

import (
	"context"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/google/uuid"
)

type RepositoryInterface interface {
	AddUser(ctx context.Context, user models.User) (*uuid.UUID, error)
}
