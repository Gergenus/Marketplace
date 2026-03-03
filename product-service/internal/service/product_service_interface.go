package service

import (
	"context"

	"github.com/Gergenus/commerce/product-service/internal/models"
)

type ServiceInterface interface {
	AddCategory(ctx context.Context, category string) (int, error)
	CreateProduct(ctx context.Context, product models.Product) (int, error)
	GetStockByID(ctx context.Context, product_id, seller_id int) (int, error)
	AddStockByID(ctx context.Context, seller_id, product_id, number int) (int, error)
	GetProductByID(ctx context.Context, id int) (models.Product, error)
}
