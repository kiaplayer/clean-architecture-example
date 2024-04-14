package dto

import (
	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/reference"
)

func SaleOrderDtoToSaleOrder(saleOrderDTO SaleOrder) *document.SaleOrder {
	saleOrder := document.SaleOrder{
		Document: document.Document{},
		Customer: reference.Customer{
			Reference: reference.Reference{
				ID: saleOrderDTO.CustomerID,
			},
		},
	}
	for _, product := range saleOrderDTO.Products {
		saleOrder.Products = append(saleOrder.Products, document.SaleOrderProduct{
			Product: reference.Product{
				Reference: reference.Reference{
					ID: product.ProductID,
				},
			},
			Quantity: int(product.Quantity),
		})
	}
	return &saleOrder
}
