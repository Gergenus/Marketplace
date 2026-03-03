package service

import "context"

type ServiceInterface interface {
	AddCategory(ctx context.Context, category string) (int, error)
}
