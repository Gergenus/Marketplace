package repository

import (
	"context"

	"github.com/Gergenus/commerce/product-service/internal/models"
)

type RepositoryInterface interface {
	AddCategory(ctx context.Context, category string) (int, error)
	GetCategoryID(ctx context.Context, category string) (int, error)
	DeleteCategoryByID(ctx context.Context, id int) error
	AllProducts(ctx context.Context) ([]*models.Product, error)
	GetStockByID(ctx context.Context, product_id int) (int, error)                           // must have
	AddStockByID(ctx context.Context, seller_id string, product_id, number int) (int, error) // must have
	ReduceStock(ctx context.Context, seller_id string, product_id, number int) (int, error)
	CreateProduct(ctx context.Context, product models.Product) (int, error) // must have
	DeleteProduct(ctx context.Context, id int) error
	UpdateProduct(ctx context.Context, product models.Product, product_id int) error
	GetProductByID(ctx context.Context, id int) (models.Product, error) // must have
	GetProductsByCategory(ctx context.Context, category string) ([]models.Product, error)
	GetProductsBySellerID(ctx context.Context, seller_id string) ([]models.Product, error)
	CheckProductExists(ctx context.Context, seller_id string, product_name string) (bool, error)
	ReserveProducts(ctx context.Context, products []models.ProductsToReserve) error
}
