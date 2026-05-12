package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Gergenus/commerce/order-service/internal/models"
	"github.com/Gergenus/commerce/order-service/internal/repository"
	"github.com/Gergenus/commerce/order-service/proto"
	"github.com/google/uuid"
)

var (
	ErrNoCartFound      = errors.New("no cart found")
	ErrOrderNotReserved = errors.New("order not reserved")
)

type OrderService struct {
	repo          repository.OrderRepoInterface
	log           *slog.Logger
	cartClient    proto.OrderServiceClient
	productClient proto.OrderServiceClient
}

func NewOrderService(repo repository.OrderRepoInterface, log *slog.Logger, cartClient proto.OrderServiceClient, productClient proto.OrderServiceClient) *OrderService {
	return &OrderService{
		repo:          repo,
		log:           log,
		cartClient:    cartClient,
		productClient: productClient,
	}
}

// returns orderId
func (o *OrderService) CreateOrder(ctx context.Context, userId uuid.UUID, deliveryAddress string) (int, error) {
	const op = "service.CreateOrder"
	log := o.log.With(slog.String("op", op))
	// gRPC req to cart service
	cartResp, err := o.cartClient.GetCart(ctx, &proto.GetCartRequest{UserId: userId.String()})
	if err != nil {
		log.Error("failed to call gRPC req to cart-service", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if !cartResp.Availablility {
		log.Warn("no cart was found")
		return 0, fmt.Errorf("%s: %w", op, ErrNoCartFound)
	}

	products := cartResp.GetOrderProducts()
	// gRPC req to product-service, it returns a slice of products with its seller_id
	ReservedResponse, err := o.productClient.ReserveOrder(ctx, &proto.ReserveOrderRequest{OrderProducts: products})
	if err != nil {
		log.Error("failed to call gRPC req to product-service", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	orderId, err := o.repo.CreateOrder(ctx, userId, float64(ReservedResponse.Price), deliveryAddress)
	if err != nil {
		log.Error("failed to create order", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	convertedProducts := convertProducts(ReservedResponse.GetProductsSeller())
	// create compensating transactions
	err = o.repo.FillOrder(ctx, orderId, convertedProducts)
	if err != nil {
		if err := o.repo.DeleteOrder(ctx, orderId); err != nil {
			log.Error("failed to compensate creating order", slog.String("error", err.Error()))
		}
		log.Error("failed to fill order", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return orderId, nil
}

// for sellers
func (o *OrderService) Orders(ctx context.Context, sellerId uuid.UUID) ([]models.OrderProduct, error) {
	const op = "service.Orders"
	log := o.log.With(slog.String("op", op))
	log.Info("getting orders")
	products, err := o.repo.Orders(ctx, sellerId)
	if err != nil {
		log.Error("failed to get seller orders", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return products, nil
}

func convertProducts(oldProducts []*proto.ProductSeller) []models.OrderProduct {
	products := []models.OrderProduct{}
	for _, d := range oldProducts {
		product := models.OrderProduct{
			ID:       int(d.GetProductId()),
			Stock:    int(d.GetStock()),
			SellerID: uuid.MustParse(d.SellerId),
		}
		products = append(products, product)
	}
	return products
}
