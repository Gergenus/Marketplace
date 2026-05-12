package repository

import "context"

type RepositoryInterface interface {
	AddToCart(ctx context.Context, UUID, productID string, stock int) error
	DeleteFromCart(ctx context.Context, UUID, productID string) error
	UpdateStock(ctx context.Context, UUID, productID string, stock int) error
	GetCart(ctx context.Context, UUID string) (map[string]string, error)
}
