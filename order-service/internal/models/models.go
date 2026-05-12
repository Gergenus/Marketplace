package models

import "github.com/google/uuid"

type OrderProduct struct {
	ID              int       `json:"id,omitempty"`
	Stock           int       `json:"stock,omitempty"`
	SellerID        uuid.UUID `json:"seller_id,omitempty"`
	DeliveryAddress string    `json:"delivery_address,omitempty"`
}
