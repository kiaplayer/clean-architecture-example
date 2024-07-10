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

	saleOrder := &document.SaleOrder{
		Customer: reference.Customer{
			Reference: reference.Reference{
				ID: 1,
			},
		},
		Products: []document.SaleOrderProduct{
			{
				Product: reference.Product{
					Reference: reference.Reference{
						ID: 1,
					},
				},
				Quantity: 1,
			},
		},
	}

	useCaseMock.EXPECT().
		Handle(ctx, saleOrder).
		Return(saleOrder, nil)

	transactorMock.EXPECT().
		RunInTx(ctx, gomock.Any()).
		DoAndReturn(
			func(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
				return fn(ctx)
			},
		)

	bodyReader := bytes.NewReader([]byte(`{"customer_id": 1, "products": [{"product_id": 1, "quantity": 1}]}`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodPost, "", bodyReader)

	// act
	handler.Handle(response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, fmt.Sprintf("SaleOrder ID = %d", saleOrder.ID), response.Body.String())
}

func TestHandle_checkAccessError(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)

	useCaseMock := mocks.NewMockuseCase(ctrl)
	transactorMock := mocks.NewMocktransactor(ctrl)
	handler := NewHandler(useCaseMock, transactorMock)

	bodyReader := bytes.NewReader([]byte(`{}`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodDelete, "", bodyReader)

	// act
	handler.Handle(response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusForbidden, response.Code)
}

func TestHandle_validateError_emptyRequest(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)

	useCaseMock := mocks.NewMockuseCase(ctrl)
	transactorMock := mocks.NewMocktransactor(ctrl)
	handler := NewHandler(useCaseMock, transactorMock)

	bodyReader := bytes.NewReader([]byte(`{}`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodPost, "", bodyReader)

	// act
	handler.Handle(response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestHandle_validateError_zeroProductID(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)

	useCaseMock := mocks.NewMockuseCase(ctrl)
	transactorMock := mocks.NewMocktransactor(ctrl)
	handler := NewHandler(useCaseMock, transactorMock)

	bodyReader := bytes.NewReader([]byte(`{"customer_id": 1, "products": [{"product_id": 0, "quantity": 1}]}`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodPost, "", bodyReader)

	// act
	handler.Handle(response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestHandle_validateError_invalidJSON(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)

	useCaseMock := mocks.NewMockuseCase(ctrl)
	transactorMock := mocks.NewMocktransactor(ctrl)
	handler := NewHandler(useCaseMock, transactorMock)

	bodyReader := bytes.NewReader([]byte(`invalid_json`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodPost, "", bodyReader)

	// act
	handler.Handle(response, request)

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

	saleOrder := &document.SaleOrder{
		Customer: reference.Customer{
			Reference: reference.Reference{
				ID: 1,
			},
		},
		Products: []document.SaleOrderProduct{
			{
				Product: reference.Product{
					Reference: reference.Reference{
						ID: 1,
					},
				},
				Quantity: 1,
			},
		},
	}

	createErr := errors.New("some error while creating sale order")

	useCaseMock.EXPECT().
		Handle(ctx, saleOrder).
		Return(nil, createErr)

	transactorMock.EXPECT().
		RunInTx(ctx, gomock.Any()).
		DoAndReturn(
			func(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
				return fn(ctx)
			},
		)

	bodyReader := bytes.NewReader([]byte(`{"customer_id": 1, "products": [{"product_id": 1, "quantity": 1}]}`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodPost, "", bodyReader)

	// act
	handler.Handle(response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
}
