package service

import (
	"context"

	"github.com/Gergenus/commerce/order-service/internal/models"
	"github.com/google/uuid"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, userId uuid.UUID, deliveryAddress string) (int, error)
	Orders(ctx context.Context, sellerId uuid.UUID) ([]models.OrderProduct, error) // for sellers
}
