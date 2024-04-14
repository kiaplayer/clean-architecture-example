package dto

type SaleOrder struct {
	CustomerID uint64             `json:"customer_id"`
	Products   []SaleOrderProduct `json:"products"`
}

type SaleOrderProduct struct {
	ProductID uint64 `json:"product_id"`
	Quantity  uint64 `json:"quantity"`
}
