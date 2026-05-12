package service

import "context"

type CartServiceInterface interface {
	AddToCart(ctx context.Context, UUID, productID string, stock int) error
	DeleteFromCart(ctx context.Context, UUID, productID string) error
	UpdateStock(ctx context.Context, UUID, productID string, stock int) error
	Cart(ctx context.Context, UUID string) (map[string]string, error)
}
