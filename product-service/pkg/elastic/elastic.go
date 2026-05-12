package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/Gergenus/commerce/product-service/internal/models"
	"github.com/Gergenus/commerce/product-service/internal/repository"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
)

const ProductIndex = "products"

type ElasticClient struct {
	ElClient *elasticsearch.Client
	log      *slog.Logger
	repo     repository.RepositoryInterface
}

func NewElasticClient(Addresses []string, username, password, crt string, log *slog.Logger, repo repository.RepositoryInterface) ElasticClient {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: Addresses,
		Username:  username,
		Password:  password,
		CACert:    []byte(crt),
	})
	if err != nil {
		panic(err)
	}
	_, err = client.Ping()
	if err != nil {
		panic(err)
	}

	return ElasticClient{ElClient: client, log: log, repo: repo}
}

func (e *ElasticClient) InitIndexation(ctx context.Context) bool {
	resp, err := esapi.IndicesExistsRequest{
		Index: []string{ProductIndex},
	}.Do(ctx, e.ElClient)

	if err != nil || resp.IsError() {
		e.ElClient.Indices.Create(ProductIndex)
		e.log.Info("index created", slog.String("index", ProductIndex))
		return true
	}
	return false
}

func (e *ElasticClient) IndexAllProducts(ctx context.Context) error {
	const op = "elastic.IndexAllProducts"
	log := e.log.With(slog.String("op", op))
	log.Info("indexing all products")
	products, err := e.repo.AllProducts(ctx)
	if err != nil {
		log.Error("failed to get products", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	for _, product := range products {
		data, err := json.Marshal(product)
		if err != nil {
			log.Error("failed to marshal product", slog.String("error", err.Error()))
			return fmt.Errorf("%s: %w", op, err)
		}
		resp, err := esapi.IndexRequest{
			Index:      ProductIndex,
			DocumentID: strconv.Itoa(product.ID),
			Body:       bytes.NewReader(data),
			Refresh:    "true",
		}.Do(ctx, e.ElClient)
		if err != nil {
			log.Error("failed to index product", slog.String("error", err.Error()), slog.Int("id", product.ID))
			return fmt.Errorf("%s: %w", op, err)
		}
		defer resp.Body.Close()
	}
	log.Info("succesful indexation")
	return nil
}

func (e *ElasticClient) IndexProduct(ctx context.Context, product models.Product) error {
	const op = "elastic.IndexProduct"
	log := e.log.With(slog.String("op", op))
	log.Info("indexing proudct")

	data, err := json.Marshal(product)
	if err != nil {
		log.Error("failed to marshal product", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	resp, err := esapi.IndexRequest{
		Index:      ProductIndex,
		DocumentID: strconv.Itoa(product.ID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}.Do(ctx, e.ElClient)
	if err != nil {
		log.Error("failed to index product", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()
	return nil
}
