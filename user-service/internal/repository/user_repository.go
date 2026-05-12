package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gergenus/commerce/user-service/internal/models"
	dbpkg "github.com/Gergenus/commerce/user-service/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrNoSessionFound     = errors.New("no session found")
	ErrVerificationFailed = errors.New("verification failed")
)

type PostgresRepository struct {
	db dbpkg.PostgresDB
}

func NewPostgresRepository(db dbpkg.PostgresDB) PostgresRepository {
	return PostgresRepository{
		db: db,
	}
}

func (p *PostgresRepository) AddUser(ctx context.Context, user models.User) (*uuid.UUID, error) {
	const op = "repository.AddUser"
	tx, err := p.db.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	uid := uuid.New()
	_, err = p.db.DB.Exec(ctx, "INSERT INTO users (id, username, email, role, hashpassword) VALUES($1, $2, $3, $4, $5)", uid.String(),
		user.Username, user.Email, user.Role, user.Password)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.Code == "23505" {
				return nil, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
			}
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &uid, nil
}

func (p *PostgresRepository) GetUser(ctx context.Context, email string) (*models.User, error) {
	const op = "repository.GetUser"
	var user models.User
	err := p.db.DB.QueryRow(ctx, "SELECT id, username, email, verified, role, hashpassword FROM users WHERE email = $1", email).Scan(&user.ID,
		&user.Username, &user.Email, &user.Verified, &user.Role, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil

}

func (p *PostgresRepository) GetUserByUUID(ctx context.Context, uuid string) (*models.User, error) {
	const op = "repository.GetUser"
	var user models.User
	err := p.db.DB.QueryRow(ctx, "SELECT id, username, email, verified, role, hashpassword FROM users WHERE id = $1", uuid).Scan(&user.ID,
		&user.Username, &user.Email, &user.Verified, &user.Role, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil

}

func (p *PostgresRepository) CreateJWTSession(ctx context.Context, userId string, refreshToken, fingerprint, ip string, expiresIn int64) error {
	const op = "repository.CreateJWTSession"
	_, err := p.db.DB.Exec(ctx, "INSERT INTO refreshsessions (userId, refreshToken, fingerprint, ip, expiresIn) VALUES($1, $2, $3, $4, $5)",
		userId, refreshToken, fingerprint, ip, expiresIn)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (p *PostgresRepository) GetRefreshSession(ctx context.Context, oldUuid string) (*models.RefreshSession, error) {
	const op = "repository.GetRefreshSession"
	tx, err := p.db.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	var session models.RefreshSession
	err = tx.QueryRow(ctx, "SELECT * FROM refreshSessions WHERE refreshToken = $1", oldUuid).Scan(&session.ID, &session.UserID, &session.ResfreshToken, &session.Fingerprint, &session.IP, &session.ExpiresIn,
		&session.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = tx.Exec(ctx, "DELETE FROM refreshSessions WHERE refreshToken = $1", oldUuid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &session, nil
}

func (p *PostgresRepository) DeleteSession(ctx context.Context, refreshToken string) error {
	const op = "repository.DeleteSession"
	row, err := p.db.DB.Exec(ctx, "DELETE FROM refreshSessions WHERE refreshToken = $1", refreshToken)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if row.RowsAffected() == 0 {
		return ErrNoSessionFound
	}
	return nil
}

func (p *PostgresRepository) Verification(ctx context.Context, email string) error {
	const op = "repository.Verification"
	fmt.Println(email, "SVOSVOSVOVSOSVo")
	row, err := p.db.DB.Exec(ctx, "UPDATE users SET verified = true WHERE email = $1", email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if row.RowsAffected() == 0 {
		return ErrVerificationFailed
	}
	return nil
}
