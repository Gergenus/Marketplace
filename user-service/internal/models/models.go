package models

type Seller struct {
	ID       int
	Name     string
	Location string
	License  string
}

type Product struct {
	ID          int
	ProductName string
	Price       float64
	SellerID    int
	CategoryID  int
}
