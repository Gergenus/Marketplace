package models

import "github.com/google/uuid"

type Product struct {
	ID          int     `json:"id,omitempty"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	SellerID    string  `json:"seller_id,omitempty"`
	CategoryID  int     `json:"category_id"`
}

type ProductsToReserve struct {
	ID       int `json:"id,omitempty"`
	Stock    int
	SellerID uuid.UUID
	Price    float64
}

type Category struct {
	Category string `json:"category"`
}

type AddStockRequest struct {
	ProductID int `json:"product_id"`
	Number    int `json:"number"`
}
