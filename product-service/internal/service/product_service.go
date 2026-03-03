package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Gergenus/commerce/product-service/internal/models"
	"github.com/Gergenus/commerce/product-service/internal/repository"
)

var (
	ErrFailedCreateProduct        = errors.New("failed to create new product")
	ErrMoreThanOneProductInstance = errors.New("more than one instance of a product")
	ErrStockNotFound              = errors.New("stock not found")
	ErrProductNotFound            = errors.New("product not found")
	ErrNoSuchCategoryExists       = errors.New("no such category exists")
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
	const op = "service.AddCategory"
	id, err := p.repo.AddCategory(ctx, category)
	if err != nil {
		p.log.Error("adding category error", slog.String("category", category))
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	p.log.Info("the category was added", slog.String("category", category))
	return id, nil
}

func (p *ProductService) CreateProduct(ctx context.Context, product models.Product) (int, error) {
	const op = "service.CreateProduct"
	p.log.Info("creating product", slog.String("product_name", product.ProductName), slog.Int("seller_id", product.SellerID),
		slog.Float64("price", product.Price), slog.Int("category_id", product.CategoryID))

	check, err := p.repo.CheckProductExists(ctx, product.SellerID, product.ProductName)
	if err != nil {
		p.log.Error("creating product error", slog.String("error", err.Error()))
		return -1, fmt.Errorf("%s: %w", op, ErrFailedCreateProduct)
	}
	if !check {
		return -1, fmt.Errorf("%s: %w", op, ErrMoreThanOneProductInstance)
	}
	id, err := p.repo.CreateProduct(ctx, product)
	if err != nil {
		p.log.Error("creating product error", slog.String("error", err.Error()))
		if errors.Is(err, repository.ErrNoSuchCategoryExists) {
			p.log.Error("creating product error", slog.String("error", err.Error()))
			return -1, fmt.Errorf("%s: %w", op, ErrNoSuchCategoryExists)
		}
		return -1, fmt.Errorf("%s: %w", op, ErrFailedCreateProduct)
	}
	return id, nil
}

func (p *ProductService) GetStockByID(ctx context.Context, product_id, seller_id int) (int, error) {
	const op = "service.GetStockByID"
	p.log.Info("getting stock", slog.Int("product_id", product_id), slog.Int("seller_id", seller_id))
	stock, err := p.repo.GetStockByID(ctx, product_id, seller_id)
	if err != nil {
		if errors.Is(err, repository.ErrStockNotFound) {
			p.log.Error("stock not found", slog.Int("product_id", product_id), slog.Int("seller_id", seller_id))
			return -1, fmt.Errorf("%s: %w", op, ErrStockNotFound)
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return stock, nil
}

func (p *ProductService) AddStockByID(ctx context.Context, seller_id, product_id, number int) (int, error) {
	const op = "service.AddStockByID"
	p.log.Info("adding stock", slog.Int("id", product_id), slog.Int("seller_id", seller_id))
	id, err := p.repo.AddStockByID(ctx, seller_id, product_id, number)
	if err != nil {
		p.log.Error("error adding stock by id", slog.String("error", err.Error()))
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (p *ProductService) GetProductByID(ctx context.Context, id int) (models.Product, error) {
	const op = "service.GetProductByID"
	p.log.Info("getting product", slog.Int("product_id", id))
	product, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoSuchProductExists) {
			p.log.Error("product not found", slog.Int("product_id", id))
			return models.Product{}, fmt.Errorf("%s: %w", op, ErrProductNotFound)
		}
		p.log.Error("error getting product", slog.String("error", err.Error()))
		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}
	return product, nil
}
