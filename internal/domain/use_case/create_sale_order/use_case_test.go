package create_sale_order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/reference"
	mocks "github.com/kiaplayer/clean-architecture-example/internal/domain/use_case/create_sale_order/mocks"
)

func TestHandle_Success(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	timeGeneratorMock := mocks.NewMocktimeGenerator(ctrl)
	numberGeneratorMock := mocks.NewMocknumberGenerator(ctrl)
	saleOrderServiceMock := mocks.NewMocksaleOrderService(ctrl)

	useCase := NewUseCase(timeGeneratorMock, numberGeneratorMock, saleOrderServiceMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Date:       time.Now().Truncate(time.Second),
			Number:     "2000-01-01-11-001",
			Company:    reference.Company{},
			AppendUser: nil,
		},
		Customer: reference.Customer{},
		Products: []document.SaleOrderProduct{
			{
				Product: reference.Product{
					Reference: reference.Reference{
						ID:     1,
						Name:   "Товар 1",
						Status: reference.StatusActive,
					},
				},
				Quantity: 1,
				Price:    100,
			},
		},
	}

	timeGeneratorMock.EXPECT().
		NowDate().
		Return(saleOrder.Date)

	numberGeneratorMock.EXPECT().
		GenerateNumber(saleOrder.Date, saleOrder.Company).
		Return(saleOrder.Number)

	saleOrderServiceMock.EXPECT().
		CreateOrder(ctx, saleOrder).
		Return(saleOrder, nil)

	// act
	actualSaleOrder, actualErr := useCase.Handle(ctx, saleOrder)

	// assert
	assert.NoError(t, actualErr)
	assert.Equal(t, saleOrder, actualSaleOrder)
}

func TestHandle_Error(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	timeGeneratorMock := mocks.NewMocktimeGenerator(ctrl)
	numberGeneratorMock := mocks.NewMocknumberGenerator(ctrl)
	saleOrderServiceMock := mocks.NewMocksaleOrderService(ctrl)

	useCase := NewUseCase(timeGeneratorMock, numberGeneratorMock, saleOrderServiceMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Date:   time.Now().Truncate(time.Second),
			Number: "0001",
		},
	}

	createErr := errors.New("create error")

	timeGeneratorMock.EXPECT().
		NowDate().
		Return(saleOrder.Date)

	numberGeneratorMock.EXPECT().
		GenerateNumber(saleOrder.Date, saleOrder.Company).
		Return(saleOrder.Number)

	saleOrderServiceMock.EXPECT().
		CreateOrder(ctx, saleOrder).
		Return(nil, createErr)

	// act
	actualSaleOrder, actualErr := useCase.Handle(ctx, saleOrder)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, actualErr, createErr.Error())
}
