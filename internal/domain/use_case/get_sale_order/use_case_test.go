package get_sale_order

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	mocks "github.com/kiaplayer/clean-architecture-example/internal/domain/use_case/get_sale_order/mocks"
)

func TestHandle_Success(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	saleOrderServiceMock := mocks.NewMocksaleOrderService(ctrl)

	useCase := NewUseCase(saleOrderServiceMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			ID: 1,
		},
	}

	saleOrderServiceMock.EXPECT().
		GetOrderByID(ctx, saleOrder.ID).
		Return(saleOrder, nil)

	// act
	actualSaleOrder, actualErr := useCase.Handle(ctx, saleOrder.ID)

	// assert
	assert.NoError(t, actualErr)
	assert.Equal(t, saleOrder, actualSaleOrder)
}

func TestHandle_Error(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	saleOrderServiceMock := mocks.NewMocksaleOrderService(ctrl)

	useCase := NewUseCase(saleOrderServiceMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			ID: 1,
		},
	}

	getErr := errors.New("get error")

	saleOrderServiceMock.EXPECT().
		GetOrderByID(ctx, saleOrder.ID).
		Return(nil, getErr)

	// act
	actualSaleOrder, actualErr := useCase.Handle(ctx, saleOrder.ID)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, actualErr, getErr.Error())
}
