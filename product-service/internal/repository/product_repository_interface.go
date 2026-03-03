package repository

import (
	"context"

	"github.com/Gergenus/commerce/user-service/internal/models"
)

type RepositoryInterface interface {
	AddCategory(ctx context.Context, category string) (int, error)
	DeleteCategoryByID(ctx context.Context, id int) (int, error)
	GetProductsByCategory(ctx context.Context, category string) ([]models.Product, error)
	GetStockByID(ctx context.Context, id int) (int, error)
	AddStockByID(ctx context.Context, id int, number int) (int, error)
	ReduceStock(ctx context.Context, id int, number int) (int, error)
	CreateProduct(ctx context.Context, product models.Product) (int, error)
	DeleteProduct(ctx context.Context, id int) (int, error)
	UpdateProduct(ctx context.Context, product models.Product) (int, error)
}
