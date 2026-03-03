package repository

import (
	"context"
	"log/slog"

	"github.com/Gergenus/commerce/user-service/internal/models"
	dbpkg "github.com/Gergenus/commerce/user-service/pkg/db"
)

type PostgresRepository struct {
	db  dbpkg.PostgresDB
	log *slog.Logger
}

func NewPostgresRepository(DB dbpkg.PostgresDB, Log *slog.Logger) PostgresRepository {
	return PostgresRepository{
		db:  DB,
		log: Log,
	}
}

func (p *PostgresRepository) AddCategory(ctx context.Context, category string) (int, error) {
	var id int
	err := p.db.DB.QueryRow("INSERT INTO categories (category) VALUES($1) RETURNING id", category).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (p *PostgresRepository) DeleteCategoryByID(ctx context.Context, id int) (int, error) {
	panic("implement")
}

func (p *PostgresRepository) GetProductByID(ctx context.Context, id int) (models.Product, error) {
	panic("implement")
}

func (p *PostgresRepository) GetProductsByCategory(ctx context.Context, category string) ([]models.Product, error) {
	panic("implement")
}

func (p *PostgresRepository) GetStockByID(ctx context.Context, id int) (int, error) {
	panic("implement")
}

func (p *PostgresRepository) AddStockByID(ctx context.Context, id int, number int) (int, error) {
	panic("implement")
}

func (p *PostgresRepository) ReduceStock(ctx context.Context, id int, number int) (int, error) {
	panic("implement")
}

func (p *PostgresRepository) CreateProduct(ctx context.Context, product models.Product) (int, error) {
	panic("implement")
}

func (p *PostgresRepository) DeleteProduct(ctx context.Context, id int) (int, error) {
	panic("implement")
}

func (p *PostgresRepository) UpdateProduct(ctx context.Context, product models.Product) (int, error) {
	panic("implement")
}
