package models

type AddCartRequest struct {
	ProductId int `json:"product_id"`
	Stock     int `json:"stock"`
}

type DeleteFromCartRequest struct {
	ProductId int `json:"product_id"`
}
