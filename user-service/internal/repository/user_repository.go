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
	ErrUserAlreadyExists = errors.New("user already exists")
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
	const op = ""
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
