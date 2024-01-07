package create_sale_order

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/reference"
	mocks "github.com/kiaplayer/clean-architecture-example/internal/handlers/create_sale_order/mocks"
)

func TestHandle_Success(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	useCaseMock := mocks.NewMockuseCase(ctrl)
	transactorMock := mocks.NewMocktransactor(ctrl)
	handler := NewHandler(useCaseMock, transactorMock)

	var inputProducts []document.SaleOrderProduct
	inputCompany := reference.Company{}
	inputCustomer := reference.Customer{}

	createdSaleOrder := &document.SaleOrder{
		Document: document.Document{
			ID: 123,
		},
	}

	useCaseMock.EXPECT().
		Handle(
			ctx,
			inputProducts,
			inputCompany,
			inputCustomer,
			nil,
		).
		Return(createdSaleOrder, nil)

	transactorMock.EXPECT().
		RunInTx(ctx, gomock.Any()).
		DoAndReturn(
			func(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
				return fn(ctx)
			},
		)

	bodyReader := bytes.NewReader([]byte(`{"company": 12, "customer": 15, "products": []}`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodPost, "", bodyReader)

	// act
	handler.Handle(ctx, response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, fmt.Sprintf("SaleOrder ID = %d", createdSaleOrder.ID), response.Body.String())
}

func TestHandle_checkAccessError(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	useCaseMock := mocks.NewMockuseCase(ctrl)
	transactorMock := mocks.NewMocktransactor(ctrl)
	handler := NewHandler(useCaseMock, transactorMock)

	bodyReader := bytes.NewReader([]byte(`{}`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodDelete, "", bodyReader)

	// act
	handler.Handle(ctx, response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusForbidden, response.Code)
}

func TestHandle_validateAndPrepareError(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	useCaseMock := mocks.NewMockuseCase(ctrl)
	transactorMock := mocks.NewMocktransactor(ctrl)
	handler := NewHandler(useCaseMock, transactorMock)

	bodyReader := bytes.NewReader([]byte(`{}`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodPut, "", bodyReader)

	// act
	handler.Handle(ctx, response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestHandle_UseCaseError(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	useCaseMock := mocks.NewMockuseCase(ctrl)
	transactorMock := mocks.NewMocktransactor(ctrl)
	handler := NewHandler(useCaseMock, transactorMock)

	var inputProducts []document.SaleOrderProduct
	inputCompany := reference.Company{}
	inputCustomer := reference.Customer{}

	createErr := errors.New("some error while creating sale order")

	useCaseMock.EXPECT().
		Handle(
			ctx,
			inputProducts,
			inputCompany,
			inputCustomer,
			nil,
		).
		Return(nil, createErr)

	transactorMock.EXPECT().
		RunInTx(ctx, gomock.Any()).
		DoAndReturn(
			func(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
				return fn(ctx)
			},
		)

	bodyReader := bytes.NewReader([]byte(`{"company": 12, "customer": 15, "products": []}`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodPost, "", bodyReader)

	// act
	handler.Handle(ctx, response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
}
