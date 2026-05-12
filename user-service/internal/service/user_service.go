package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Gergenus/commerce/user-service/internal/brocker"
	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/Gergenus/commerce/user-service/internal/repository"
	"github.com/Gergenus/commerce/user-service/pkg/hash"
	"github.com/Gergenus/commerce/user-service/pkg/jwtpkg"
	"github.com/Gergenus/commerce/user-service/pkg/utils"
	"github.com/Gergenus/commerce/user-service/pkg/validation"
	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrIncorrectEmail        = errors.New("incorrect email")
	ErrPasswordMismatch      = errors.New("passwords mismatch")
	ErrTokenExpired          = errors.New("token expired")
	ErrInvalidRefreshSession = errors.New("invalid refresh session")
	ErrNoSessionFound        = errors.New("no session found")
)

type UserService struct {
	log           *slog.Logger
	repo          repository.RepositoryInterface
	jwtToken      jwtpkg.UserJWTInterface
	RefreshTTl    time.Duration
	EventProducer brocker.BrockerInterface
}

func NewUserService(log *slog.Logger, repo repository.RepositoryInterface,
	jwtToken jwtpkg.UserJWTInterface, RefreshTTl time.Duration, EventProducer brocker.BrockerInterface) UserService {
	return UserService{log: log, repo: repo, jwtToken: jwtToken, RefreshTTl: RefreshTTl, EventProducer: EventProducer}
}

func (u *UserService) AddUser(ctx context.Context, user models.User) (*uuid.UUID, error) {
	const op = "service.AddUser"
	u.log.With(slog.String("op", op))
	u.log.Info("Creating user", slog.String("email", user.Email))

	if !validation.IsEmailValid(user.Email) {
		u.log.Error("failed to validate email")
		return nil, fmt.Errorf("%s: %w", op, ErrIncorrectEmail)
	}

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

func (u *UserService) Login(ctx context.Context, email, password, userAgent, ip string) (string, string, error) {
	const op = "service.Login"
	u.log.With(slog.String("op", op))
	u.log.Info("logigng the user", slog.String("email", email))

	user, err := u.repo.GetUser(ctx, email)
	if err != nil {
		u.log.Error("failed to get user", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	if !hash.CheckPassword(user.Password, password) {
		u.log.Warn("passwords mismatch")
		return "", "", fmt.Errorf("%s: %w", op, ErrPasswordMismatch)
	}
	AccessToken, err := u.jwtToken.GenerateAccessToken(*user)
	if err != nil {
		u.log.Error("failed to create access token", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	RefreshToken := uuid.New().String()

	fingerprint := utils.CreateFingerprint(ip, userAgent)

	err = u.repo.CreateJWTSession(ctx, user.ID.String(), RefreshToken, fingerprint, ip, time.Now().Add(u.RefreshTTl).Unix())
	if err != nil {
		u.log.Error("failed to create jwt session", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	return AccessToken, RefreshToken, nil
}

// returns both new refreshtoken and accesToken
func (u *UserService) RefreshToken(ctx context.Context, oldUuid uuid.UUID, userAgent, ip string, oldAccessToken string) (*uuid.UUID, string, error) {
	const op = "service.RefreshToken"
	log := u.log.With(slog.String("op", op))
	log.Info("refreshing token")
	session, err := u.repo.GetRefreshSession(ctx, oldUuid.String())
	if err != nil {
		log.Error("getting refresh token error", slog.String("error", err.Error()))
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}
	if session.ExpiresIn < time.Now().Unix() {
		return nil, "", ErrTokenExpired
	}
	fingerprint := utils.CreateFingerprint(ip, userAgent)

	if session.Fingerprint != fingerprint {
		return nil, "", ErrInvalidRefreshSession
	}
	newRefresh := uuid.New()
	err = u.repo.CreateJWTSession(ctx, session.UserID.String(), newRefresh.String(), session.Fingerprint,
		session.IP, time.Now().Add(u.RefreshTTl).Unix())
	if err != nil {
		log.Error("creating jwt session error", slog.String("error", err.Error()))
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := u.jwtToken.RegenerateToken(oldAccessToken)
	if err != nil {
		log.Error("creating jwt token error", slog.String("error", err.Error()))
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}
	return &newRefresh, token, nil
}

func (u *UserService) Logout(ctx context.Context, refreshToken string) error {
	const op = "service.Logout"
	log := u.log.With(slog.String("op", op))
	log.Info("deleting session", slog.String("refresh", refreshToken))
	err := u.repo.DeleteSession(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, repository.ErrNoSessionFound) {
			return ErrNoSessionFound
		}
		log.Error("deleteing session error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("deleted session", slog.String("refresh", refreshToken))
	return nil
}

// the token consists of expiration time and email
func (u *UserService) Verification(ctx context.Context, token string) error {
	const op = "service.Verification"
	log := u.log.With(slog.String("op", op))
	log.Info("verificating the profile", slog.String("token", token))

	email, err := u.jwtToken.ParseMailToken(token)
	if err != nil {
		log.Error("token parsing error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	err = u.repo.Verification(ctx, email)
	if err != nil {
		log.Error("verification error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
