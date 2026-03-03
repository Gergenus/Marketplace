package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gergenus/commerce/product-service/internal/models"
	dbpkg "github.com/Gergenus/commerce/product-service/pkg/db"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresRepository struct {
	db dbpkg.PostgresDB
}

func NewPostgresRepository(DB dbpkg.PostgresDB) PostgresRepository {
	return PostgresRepository{
		db: DB,
	}
}

func (p *PostgresRepository) AddCategory(ctx context.Context, category string) (int, error) {
	var id int
	err := p.db.DB.QueryRow(ctx, "INSERT INTO categories (category) VALUES($1) RETURNING id", category).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (p *PostgresRepository) DeleteCategoryByID(ctx context.Context, id int) error {
	const op = "repository.DeleteCategoryByID"
	_, err := p.db.DB.Exec(ctx, "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
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
	const op = "repository.CreateProduct"
	var id int
	err := p.db.DB.QueryRow(ctx, "INSERT INTO product_list (product_name, price, seller_id, category_id) VALUES($1, $2, $3, $4) RETURNING id",
		product.ProductName, product.Price, product.SellerID, product.CategoryID).Scan(&id)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			fmt.Println(pgxErr.Code)
			fmt.Println(pgxErr.Message)
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (p *PostgresRepository) DeleteProduct(ctx context.Context, id int) error {
	const op = "repository.DeleteProduct"
	_, err := p.db.DB.Exec(ctx, "DELETE FROM product_list WHERE id = $1", id)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			fmt.Println(pgxErr.Code)
			fmt.Println(pgxErr.Message)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (p *PostgresRepository) UpdateProduct(ctx context.Context, product models.Product) error {
	const op = "repository.UpdateProduct"
	_, err := p.db.DB.Exec(ctx, "UPDATE product_list SET product_name = $1, price = $2, category_id = $3 RETURNING id",
		product.ProductName, product.Price, product.CategoryID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PostgresRepository) GetProductByID(ctx context.Context, id int) (models.Product, error) {
	panic("implement")
}

func (p *PostgresRepository) GetProductsByCategory(ctx context.Context, category string) ([]models.Product, error) {
	panic("implement")
}
