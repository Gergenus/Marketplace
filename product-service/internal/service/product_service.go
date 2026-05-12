package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/Gergenus/commerce/product-service/internal/models"
	"github.com/Gergenus/commerce/product-service/internal/repository"
	"github.com/Gergenus/commerce/product-service/pkg/elastic"
	"github.com/google/uuid"
)

var (
	ErrFailedCreateProduct        = errors.New("failed to create new product")
	ErrMoreThanOneProductInstance = errors.New("more than one instance of a product")
	ErrStockNotFound              = errors.New("stock not found")
	ErrProductNotFound            = errors.New("product not found")
	ErrNoSuchCategoryExists       = errors.New("no such category exists")
	ErrCategoryAlreadyExists      = errors.New("category already exists")
)

type ProductService struct {
	log     *slog.Logger
	repo    repository.RepositoryInterface
	eClient *elastic.ElasticClient
}

func NewProductService(log *slog.Logger, repo repository.RepositoryInterface, eClient *elastic.ElasticClient) ProductService {
	return ProductService{
		log:     log,
		repo:    repo,
		eClient: eClient,
	}
}

func (p *ProductService) AddCategory(ctx context.Context, category string) (int, error) {
	p.log.Info("adding category", slog.String("category", category))
	const op = "service.AddCategory"
	id, err := p.repo.AddCategory(ctx, category)
	if err != nil {
		p.log.Error("adding category error", slog.String("category", category), slog.String("error", err.Error()))
		return -1, fmt.Errorf("%s: %w", op, ErrCategoryAlreadyExists)
	}
	p.log.Info("the category was added", slog.String("category", category))
	return id, nil
}

func (p *ProductService) CreateProduct(ctx context.Context, product models.Product) (int, error) {
	const op = "service.CreateProduct"
	p.log.Info("creating product", slog.String("product_name", product.ProductName), slog.String("seller_id", product.SellerID),
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
	product.ID = id
	err = p.eClient.IndexProduct(ctx, product)
	if err != nil {
		p.log.Error("indexing product error", slog.String("error", err.Error()))
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (p *ProductService) GetStockByID(ctx context.Context, product_id int) (int, error) {
	const op = "service.GetStockByID"
	p.log.Info("getting stock", slog.Int("product_id", product_id))
	stock, err := p.repo.GetStockByID(ctx, product_id)
	if err != nil {
		if errors.Is(err, repository.ErrStockNotFound) {
			p.log.Error("stock not found", slog.Int("product_id", product_id))
			return -1, fmt.Errorf("%s: %w", op, ErrStockNotFound)
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return stock, nil
}

func (p *ProductService) AddStockByID(ctx context.Context, seller_id string, product_id, number int) (int, error) {
	const op = "service.AddStockByID"
	p.log.Info("adding stock", slog.Int("id", product_id), slog.String("seller_id", seller_id))
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

func (p *ProductService) ReserveProducts(ctx context.Context, products []models.ProductsToReserve) ([]models.ProductsToReserve, error) {
	const op = "service.ReserveOrderReserveOrder"
	log := p.log.With(slog.String("op", op))
	log.Info("reserving products")
	err := p.repo.ReserveProducts(ctx, products)
	if err != nil {
		log.Error("failed to reserve products")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for index := range products {
		product, err := p.repo.GetProductByID(ctx, products[index].ID)
		if err != nil {
			log.Error("failed to get seller_id")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		sellerID, err := uuid.Parse(product.SellerID)
		if err != nil {
			log.Error("failed to parse uuid", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		products[index].SellerID = sellerID
		products[index].Price = product.Price
	}
	log.Info("succesfully reserved products")
	return products, nil
}

func (p *ProductService) Products(ctx context.Context, name string, offset, limit string) ([]models.Product, error) {
	const op = "service.Products"
	log := p.log.With(slog.String("op", op))
	log.Info("searching for products")
	from, err := strconv.Atoi(offset)
	if err != nil {
		log.Error("failed to encode query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	q := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     name,
				"fields":    []string{"product_name"},
				"fuzziness": "AUTO",
				"operator":  "or",
			},
		},
		"size": limit,
		"from": strconv.Itoa(from - 1),
		"sort": []string{},
	}

	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(q)
	if err != nil {
		log.Error("failed to encode query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := p.eClient.ElClient.Search(
		p.eClient.ElClient.Search.WithIndex(elastic.ProductIndex),
		p.eClient.ElClient.Search.WithBody(&buf),
	)
	defer resp.Body.Close()
	if err != nil {
		log.Error("failed to search", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if resp.IsError() {
		log.Error("failed to search")
	}
	var r map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Error("failed to decode reply")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	response := []models.Product{}

	if hits, ok := r["hits"].(map[string]interface{}); ok {
		if wrappedHits, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range wrappedHits {
				prod := models.Product{}
				if hitMap, ok := hit.(map[string]interface{}); ok {
					if data, ok := hitMap["_source"].(map[string]interface{}); ok {
						prod.ID = int(data["id"].(float64))
						prod.CategoryID = int(data["category_id"].(float64))
						prod.Price = data["price"].(float64)
						prod.ProductName = data["product_name"].(string)
						prod.SellerID = data["seller_id"].(string)
					}
				}
				response = append(response, prod)
			}
		}
	}
	return response, nil
}
