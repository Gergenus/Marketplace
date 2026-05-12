package service

import (
	"context"

	"github.com/Gergenus/commerce/product-service/internal/models"
)

type ServiceInterface interface {
	AddCategory(ctx context.Context, category string) (int, error)
	CreateProduct(ctx context.Context, product models.Product) (int, error)
	GetStockByID(ctx context.Context, product_id int) (int, error)
	AddStockByID(ctx context.Context, seller_id string, product_id, number int) (int, error)
	GetProductByID(ctx context.Context, id int) (models.Product, error)
	ReserveProducts(ctx context.Context, products []models.ProductsToReserve) ([]models.ProductsToReserve, error)
	Products(ctx context.Context, name string, offset, limit string) ([]models.Product, error)
}
