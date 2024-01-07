package sale_order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	mocks "github.com/kiaplayer/clean-architecture-example/internal/domain/service/sale_order/mocks"
)

func TestCreateOrder_Success(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	repositoryMock := mocks.NewMockrepository(ctrl)

	service := NewService(repositoryMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Date:   time.Now().Truncate(time.Second),
			Number: "0001",
		},
	}

	repositoryMock.EXPECT().
		CreateOrder(ctx, saleOrder).
		Return(saleOrder, nil)

	// act
	actualSaleOrder, actualErr := service.CreateOrder(ctx, saleOrder)

	// assert
	assert.NoError(t, actualErr)
	assert.Equal(t, saleOrder, actualSaleOrder)
}

func TestCreateOrder_ValidateError(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	repositoryMock := mocks.NewMockrepository(ctrl)

	service := NewService(repositoryMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Date:   time.Now().Truncate(time.Second),
			Number: "0001",
			Status: 999,
		},
	}

	// act
	actualSaleOrder, actualErr := service.CreateOrder(ctx, saleOrder)

	// assert
	assert.Equal(t, saleOrder, actualSaleOrder)
	assert.ErrorContains(t, actualErr, "bad status")
}

func TestCreateOrder_CreateError(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	repositoryMock := mocks.NewMockrepository(ctrl)

	service := NewService(repositoryMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Date:   time.Now().Truncate(time.Second),
			Number: "0001",
		},
	}

	createErr := errors.New("create error")

	repositoryMock.EXPECT().
		CreateOrder(ctx, saleOrder).
		Return(nil, createErr)

	// act
	actualSaleOrder, actualErr := service.CreateOrder(ctx, saleOrder)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, actualErr, createErr.Error())
}
