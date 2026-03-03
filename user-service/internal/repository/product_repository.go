package repository

import (
	"context"
	"log/slog"

	dbpkg "github.com/Gergenus/commerce/user-service/pkg/db"
)

type PostgresRepository struct {
	DB  dbpkg.PostgresDB
	Log *slog.Logger
}

func NewPostgresRepository(DB dbpkg.PostgresDB, Log *slog.Logger) PostgresRepository {
	return PostgresRepository{
		DB:  DB,
		Log: Log,
	}
}

func (p *PostgresRepository) AddCategory(ctx context.Context) error {

}

func (p *PostgresRepository) DeleteCategory(ctx context.Context) error {

}

func (p *PostgresRepository) AddSeller(ctx context.Context) error {

}

func (p *PostgresRepository) DeleteSeller(ctx context.Context) error {

}

func (p *PostgresRepository) UpdateSeller(ctx context.Context) error {

}

func (p *PostgresRepository) GetProductByID(ctx context.Context) error {

}

func (p *PostgresRepository) GetProductsByCategory(ctx context.Context) error {

}

func (p *PostgresRepository) GetStockByID(ctx context.Context) error {

}

func (p *PostgresRepository) AddStock(ctx context.Context) error {

}

func (p *PostgresRepository) ReduceStock(ctx context.Context) error {

}

func (p *PostgresRepository) CreateProduct(ctx context.Context) error {

}

func (p *PostgresRepository) DeleteProduct(ctx context.Context) error {

}

func (p *PostgresRepository) UpdateProduct(ctx context.Context) error {

}
