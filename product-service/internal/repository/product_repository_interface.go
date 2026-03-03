package repository

import (
	"context"

	"github.com/Gergenus/commerce/product-service/internal/models"
)

type RepositoryInterface interface {
	AddCategory(ctx context.Context, category string) (int, error)
	GetCategoryID(ctx context.Context, category string) (int, error)
	DeleteCategoryByID(ctx context.Context, id int) error
	GetStockByID(ctx context.Context, product_id, seller_id int) (int, error)
	AddStockByID(ctx context.Context, seller_id, product_id, number int) (int, error)
	ReduceStock(ctx context.Context, seller_id, product_id, number int) (int, error)
	CreateProduct(ctx context.Context, product models.Product) (int, error)
	DeleteProduct(ctx context.Context, id int) error
	UpdateProduct(ctx context.Context, product models.Product, product_id int) error
	GetProductByID(ctx context.Context, id int) (models.Product, error)
	GetProductsByCategory(ctx context.Context, category string) ([]models.Product, error)
	GetProductsBySellerID(ctx context.Context, seller_id int) ([]models.Product, error)
}
