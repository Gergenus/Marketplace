package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Gergenus/commerce/user-service/internal/repository"
)

type ProductService struct {
	log  *slog.Logger
	repo repository.RepositoryInterface
}

func NewProductService(log *slog.Logger, repo repository.RepositoryInterface) ProductService {
	return ProductService{
		log:  log,
		repo: repo,
	}
}

func (p *ProductService) AddCategory(ctx context.Context, category string) (int, error) {
	p.log.Info("adding category", slog.String("category", category))
	const op = "repository.AddCategory"
	id, err := p.repo.AddCategory(ctx, category)
	if err != nil {
		p.log.Error("adding category error", slog.String("category", category))
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	p.log.Info("the category was added", slog.String("category", category))
	return id, nil
}
