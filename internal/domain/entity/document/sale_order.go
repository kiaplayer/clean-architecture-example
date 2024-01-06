package document

import "github.com/kiaplayer/clean-architecture-example/internal/domain/entity/reference"

type SaleOrder struct {
	Document
	Customer reference.Customer
	Products []SaleOrderProduct
}

type SaleOrderProduct struct {
	Product  reference.Product
	Quantity int
	Price    float32
}
