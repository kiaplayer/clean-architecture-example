package sale_order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/reference"
	mocks "github.com/kiaplayer/clean-architecture-example/internal/domain/service/sale_order/mocks"
)

func TestCreateOrder_Success(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	repositoryMock := mocks.NewMockrepository(ctrl)
	productRepositoryMock := mocks.NewMockproductRepository(ctrl)

	service := NewService(repositoryMock, productRepositoryMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Date:   time.Now().Truncate(time.Second),
			Number: "0001",
		},
		Products: []document.SaleOrderProduct{
			{
				Product: reference.Product{
					Reference: reference.Reference{
						ID:     1,
						Name:   "Keyboard",
						Status: reference.StatusActive,
					},
				},
				Quantity: 1,
				Price:    150.5,
			},
		},
	}

	repositoryMock.
		EXPECT().
		CreateOrder(ctx, saleOrder).
		Return(saleOrder, nil)

	productRepositoryMock.
		EXPECT().
		Exists(ctx, saleOrder.Products[0].Product.ID).
		Return(true, nil)

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
	productRepositoryMock := mocks.NewMockproductRepository(ctrl)

	service := NewService(repositoryMock, productRepositoryMock)

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
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, actualErr, "bad status")
}

func TestCreateOrder_ValidateError_BadProductID(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	repositoryMock := mocks.NewMockrepository(ctrl)
	productRepositoryMock := mocks.NewMockproductRepository(ctrl)

	service := NewService(repositoryMock, productRepositoryMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Date:   time.Now().Truncate(time.Second),
			Number: "0001",
			Status: document.StatusDraft,
		},
		Products: []document.SaleOrderProduct{
			{
				Product: reference.Product{
					Reference: reference.Reference{
						ID:     999,
						Name:   "Keyboard",
						Status: reference.StatusActive,
					},
				},
				Quantity: 1,
				Price:    150.5,
			},
		},
	}

	productRepositoryMock.
		EXPECT().
		Exists(ctx, saleOrder.Products[0].Product.ID).
		Return(false, nil)

	// act
	actualSaleOrder, actualErr := service.CreateOrder(ctx, saleOrder)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, actualErr, "bad product id: 999")
}

func TestCreateOrder_ValidateError_CheckProductIDsError(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	repositoryMock := mocks.NewMockrepository(ctrl)
	productRepositoryMock := mocks.NewMockproductRepository(ctrl)

	service := NewService(repositoryMock, productRepositoryMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Date:   time.Now().Truncate(time.Second),
			Number: "0001",
			Status: document.StatusDraft,
		},
		Products: []document.SaleOrderProduct{
			{
				Product: reference.Product{
					Reference: reference.Reference{
						ID:     1,
						Name:   "Keyboard",
						Status: reference.StatusActive,
					},
				},
				Quantity: 1,
				Price:    150.5,
			},
		},
	}

	checkErr := errors.New("some db error")

	productRepositoryMock.
		EXPECT().
		Exists(ctx, saleOrder.Products[0].Product.ID).
		Return(false, checkErr)

	// act
	actualSaleOrder, actualErr := service.CreateOrder(ctx, saleOrder)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, actualErr, checkErr.Error())
}

func TestCreateOrder_CreateError(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	repositoryMock := mocks.NewMockrepository(ctrl)
	productRepositoryMock := mocks.NewMockproductRepository(ctrl)

	service := NewService(repositoryMock, productRepositoryMock)

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

func TestGetOrderByID_Success(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	repositoryMock := mocks.NewMockrepository(ctrl)
	productRepositoryMock := mocks.NewMockproductRepository(ctrl)

	service := NewService(repositoryMock, productRepositoryMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			ID:     1,
			Date:   time.Now().Truncate(time.Second),
			Number: "0001",
		},
	}

	repositoryMock.EXPECT().
		GetByID(ctx, saleOrder.ID).
		Return(saleOrder, nil)

	// act
	actualSaleOrder, actualErr := service.GetOrderByID(ctx, saleOrder.ID)

	// assert
	assert.NoError(t, actualErr)
	assert.Equal(t, saleOrder, actualSaleOrder)
}
