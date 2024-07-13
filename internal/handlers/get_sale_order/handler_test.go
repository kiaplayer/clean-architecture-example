package get_sale_order

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	mocks "github.com/kiaplayer/clean-architecture-example/internal/handlers/get_sale_order/mocks"
)

func TestHandle_Success(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	useCaseMock := mocks.NewMockuseCase(ctrl)
	handler := NewHandler(useCaseMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			ID: 123,
		},
	}

	useCaseMock.EXPECT().
		Handle(ctx, saleOrder.ID).
		Return(saleOrder, nil)

	bodyReader := bytes.NewReader([]byte(`{}`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodPost, fmt.Sprintf("?id=%d", saleOrder.ID), bodyReader)

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
	handler := NewHandler(useCaseMock)

	bodyReader := bytes.NewReader([]byte(`{}`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodDelete, "", bodyReader)

	// act
	handler.Handle(response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusForbidden, response.Code)
	assert.Equal(t, "access denied", response.Body.String())
}

func TestHandle_validateAndPrepareError_BadID(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)

	useCaseMock := mocks.NewMockuseCase(ctrl)
	handler := NewHandler(useCaseMock)

	bodyReader := bytes.NewReader([]byte(`{}`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodPut, "?id=bad", bodyReader)

	// act
	handler.Handle(response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "bad id")
}

func TestHandle_validateAndPrepareError_NegativeID(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)

	useCaseMock := mocks.NewMockuseCase(ctrl)
	handler := NewHandler(useCaseMock)

	bodyReader := bytes.NewReader([]byte(`{}`))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodPut, "?id=-11", bodyReader)

	// act
	handler.Handle(response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "bad id")
}

func TestHandle_UseCaseError(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	useCaseMock := mocks.NewMockuseCase(ctrl)
	handler := NewHandler(useCaseMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			ID: 123,
		},
	}

	getErr := errors.New("get errror")

	useCaseMock.EXPECT().Handle(ctx, saleOrder.ID).Return(nil, getErr)

	bodyReader := bytes.NewReader([]byte(``))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodPost, fmt.Sprintf("?id=%d", saleOrder.ID), bodyReader)

	// act
	handler.Handle(response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Empty(t, response.Body.String())
}

func TestHandle_UseCaseNotFound(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	useCaseMock := mocks.NewMockuseCase(ctrl)
	handler := NewHandler(useCaseMock)

	var saleOrderID uint64 = 123

	useCaseMock.EXPECT().Handle(ctx, saleOrderID).Return(nil, nil)

	bodyReader := bytes.NewReader([]byte(``))
	response := httptest.NewRecorder()
	request, requestErr := http.NewRequest(http.MethodPost, fmt.Sprintf("?id=%d", saleOrderID), bodyReader)

	// act
	handler.Handle(response, request)

	// assert
	assert.NoError(t, requestErr)
	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, "sale order not found", response.Body.String())
}
