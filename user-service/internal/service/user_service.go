package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/Gergenus/commerce/user-service/internal/repository"
	"github.com/Gergenus/commerce/user-service/pkg/hash"
	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserService struct {
	log  *slog.Logger
	repo repository.RepositoryInterface
}

func NewUserService(log *slog.Logger, repo repository.RepositoryInterface) UserService {
	return UserService{log: log, repo: repo}
}

func (u *UserService) AddUser(ctx context.Context, user models.User) (*uuid.UUID, error) {
	const op = "service.AddUser"
	u.log.With(slog.String("op", op))
	u.log.Info("Creating user", slog.String("email", user.Email))
	hashPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		slog.Error("failed to generate hashpassword", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	user.Password = hashPassword
	uid, err := u.repo.AddUser(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			u.log.Error("user already exists error", slog.String("email", user.Email))
			return nil, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}
		u.log.Error("creating user error", slog.String("email", user.Email), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	u.log.Info("created user", slog.String("email", user.Email))
	return uid, nil
}
