package repository

import (
	"context"

	"github.com/Gergenus/commerce/order-service/internal/models"
	"github.com/google/uuid"
)

type OrderRepoInterface interface {
	CreateOrder(ctx context.Context, userId uuid.UUID, price float64, deliveryAddress string) (int, error)
	FillOrder(ctx context.Context, orderId int, products []models.OrderProduct) error
	DeleteOrder(ctx context.Context, orderId int) error
	Orders(ctx context.Context, sellerId uuid.UUID) ([]models.OrderProduct, error)
}
