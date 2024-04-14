package dto

import (
	"reflect"
	"testing"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/reference"
)

const (
	defaultCustomerID = 1
	defaultProductID  = 2
	defaultQuantity   = 10
)

func TestSaleOrderDtoToSaleOrder(t *testing.T) {
	type args struct {
		saleOrderDTO SaleOrder
	}
	tests := []struct {
		name string
		args args
		want *document.SaleOrder
	}{
		{
			name: "with products",
			args: args{
				saleOrderDTO: SaleOrder{
					CustomerID: defaultCustomerID,
					Products: []SaleOrderProduct{
						{
							ProductID: defaultProductID,
							Quantity:  defaultQuantity,
						},
					},
				},
			},
			want: &document.SaleOrder{
				Customer: reference.Customer{
					Reference: reference.Reference{
						ID: defaultCustomerID,
					},
				},
				Products: []document.SaleOrderProduct{
					{
						Product: reference.Product{
							Reference: reference.Reference{
								ID: defaultProductID,
							},
						},
						Quantity: defaultQuantity,
						Price:    0,
					},
				},
			},
		},
		{
			name: "without products",
			args: args{
				saleOrderDTO: SaleOrder{
					CustomerID: defaultCustomerID,
				},
			},
			want: &document.SaleOrder{
				Customer: reference.Customer{
					Reference: reference.Reference{
						ID: defaultCustomerID,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SaleOrderDtoToSaleOrder(tt.args.saleOrderDTO)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaleOrderDtoToSaleOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}
