package handlers

import (
	"context"
	"strconv"

	"github.com/Gergenus/commerce/cart-service/proto"
)

func (c CartHandler) GetCart(ctx context.Context, in *proto.GetCartRequest) (*proto.GetCartResponse, error) {
	cartMap, err := c.srv.Cart(ctx, in.GetUserId())
	if err != nil {
		return &proto.GetCartResponse{}, err
	}
	if len(cartMap) == 0 {
		return &proto.GetCartResponse{Availablility: false}, nil
	}
	resp := []*proto.OrderProduct{}
	for productId, stock := range cartMap {
		productID, _ := strconv.Atoi(productId)
		Stock, _ := strconv.Atoi(stock)
		resp = append(resp, &proto.OrderProduct{ProductId: int64(productID), Stock: int64(Stock)})
	}
	return &proto.GetCartResponse{Availablility: true, OrderProducts: resp}, nil
}
