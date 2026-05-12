package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Gergenus/commerce/cart-service/internal/repository"
	"github.com/Gergenus/commerce/cart-service/proto"
)

var (
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrAlredyAdded        = errors.New("product already added")
	ErrNoCartProductFound = errors.New("no cart product found")
)

type CartService struct {
	log        *slog.Logger
	repo       repository.RepositoryInterface
	grpcClient proto.AvailablilityServiceClient
}

func NewCartService(log *slog.Logger, repo repository.RepositoryInterface, grpcClient proto.AvailablilityServiceClient) *CartService {
	return &CartService{log: log, repo: repo, grpcClient: grpcClient}
}

func (c CartService) AddToCart(ctx context.Context, UUID, productID string, stock int) error {
	const op = "service.AddToCart"
	log := c.log.With(slog.String("op", op))
	log.Info("sending grpc price grpc request", slog.String("productID", productID))
	response, err := c.grpcClient.IsAvailable(ctx, &proto.AvailablilityRequest{ProductId: productID, Stock: int64(stock)})
	if err != nil {
		log.Error("grpc request error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	if !response.Availablility {
		log.Warn("insufficient stock error", slog.String("productID", productID))
		return ErrInsufficientStock
	}
	err = c.repo.AddToCart(ctx, UUID, productID, stock)
	if err != nil {
		if errors.Is(err, repository.ErrAlredyAdded) {
			log.Warn("product alredy added", slog.String("productID", productID))
			return ErrAlredyAdded
		}
		log.Error("adding to the cart error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("item has been added to the cart", slog.String("productID", productID), slog.Int("stock", stock))
	return nil
}

func (c CartService) DeleteFromCart(ctx context.Context, UUID, productID string) error {
	const op = "service.DeleteFromCart"
	log := c.log.With(slog.String("op", op))
	log.Info("deleting from the cart", slog.String("UUID", UUID))
	err := c.repo.DeleteFromCart(ctx, UUID, productID)
	if err != nil {
		if errors.Is(err, repository.ErrNoCartProductFound) {
			log.Warn("no cart product found", slog.String("productID", productID))
			return ErrNoCartProductFound
		}
		log.Error("deleting product from the cart error", slog.String("UUID", UUID))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("product has been deleted", slog.String("UUID", UUID))
	return nil
}

func (c CartService) UpdateStock(ctx context.Context, UUID, productID string, stock int) error {
	const op = "service.UpdateStock"
	log := c.log.With(slog.String("op", op))
	log.Info("updating the cart", slog.String("UUID", UUID))
	err := c.repo.UpdateStock(ctx, UUID, productID, stock)
	if err != nil {
		log.Error("updating the cart stock error", slog.String("UUID", UUID))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (c CartService) Cart(ctx context.Context, UUID string) (map[string]string, error) {
	const op = "service.GetCart"
	log := c.log.With(slog.String("op", op))
	log.Info("getting the cart", slog.String("UUID", UUID))
	cart, err := c.repo.GetCart(ctx, UUID)
	if err != nil {
		log.Error("getting cart error", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return cart, nil
}
