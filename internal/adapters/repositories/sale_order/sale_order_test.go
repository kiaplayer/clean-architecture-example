package sale_order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mocks "github.com/kiaplayer/clean-architecture-example/internal/adapters/repositories/sale_order/mocks"
	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	"github.com/kiaplayer/clean-architecture-example/pkg/helpers"
)

func TestCreateOrder_Success(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	queryExecutorMock := mocks.NewMockQueryExecutor(ctrl)

	repository := NewRepository(queryExecutorMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Number: "123",
			Date:   time.Now(),
			Status: document.StatusDraft,
		},
	}

	lastInsertID := int64(123)

	insertResult := sqlmock.NewResult(lastInsertID, 1)

	queryExecutorMock.EXPECT().
		ExecContext(
			ctx,
			"INSERT INTO sale_order (number, date, status) VALUES (?, ?, ?)",
			saleOrder.Number,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Status,
		).
		Return(insertResult, nil)

	// act
	updatedSaleOrder, createErr := repository.CreateOrder(ctx, saleOrder)

	// assert
	assert.NoError(t, createErr)
	assert.NotNil(t, updatedSaleOrder)
	assert.Equal(t, updatedSaleOrder.ID, uint64(lastInsertID))
}

func TestCreateOrder_InsertError(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	queryExecutorMock := mocks.NewMockQueryExecutor(ctrl)

	repository := NewRepository(queryExecutorMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Number: "123",
			Date:   time.Now(),
			Status: document.StatusDraft,
		},
	}

	insertError := errors.New("insert error")

	queryExecutorMock.EXPECT().
		ExecContext(
			ctx,
			"INSERT INTO sale_order (number, date, status) VALUES (?, ?, ?)",
			saleOrder.Number,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Status,
		).
		Return(nil, insertError)

	// act
	updatedSaleOrder, createErr := repository.CreateOrder(ctx, saleOrder)

	// assert
	assert.ErrorContains(t, createErr, insertError.Error())
	assert.Nil(t, updatedSaleOrder)
}

func TestCreateOrder_LastInsertIDError(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	queryExecutorMock := mocks.NewMockQueryExecutor(ctrl)

	repository := NewRepository(queryExecutorMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Number: "123",
			Date:   time.Now(),
			Status: document.StatusDraft,
		},
	}

	lastInsertIDError := errors.New("last insert id error")

	insertResult := sqlmock.NewErrorResult(lastInsertIDError)

	queryExecutorMock.EXPECT().
		ExecContext(
			ctx,
			"INSERT INTO sale_order (number, date, status) VALUES (?, ?, ?)",
			saleOrder.Number,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Status,
		).
		Return(insertResult, nil)

	// act
	updatedSaleOrder, createErr := repository.CreateOrder(ctx, saleOrder)

	// assert
	assert.ErrorContains(t, createErr, lastInsertIDError.Error())
	assert.Nil(t, updatedSaleOrder)
}

func TestGetByID_Success(t *testing.T) {
	// arrange
	ctx := context.Background()

	db, mock, _ := sqlmock.New()

	repository := NewRepository(db)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			ID:     100,
			Number: "123",
			Date:   time.Now().Truncate(time.Second),
			Status: document.StatusDraft,
		},
	}

	rowsResult := sqlmock.NewRows([]string{
		"id",
		"date",
		"number",
		"status",
	}).AddRow(
		saleOrder.ID,
		helpers.TimeToString(saleOrder.Date),
		saleOrder.Number,
		saleOrder.Status,
	)

	mock.
		ExpectQuery("SELECT id, date, number, status FROM sale_order WHERE id = ?").
		WithArgs(saleOrder.ID).
		WillReturnRows(rowsResult)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.NoError(t, getErr)
	assert.Equal(t, saleOrder, actualSaleOrder)
}

func TestGetByID_QueryError(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	queryExecutorMock := mocks.NewMockQueryExecutor(ctrl)

	repository := NewRepository(queryExecutorMock)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			ID:     100,
			Number: "123",
			Date:   time.Now().Truncate(time.Second),
			Status: document.StatusDraft,
		},
	}

	queryError := errors.New("query error")

	queryExecutorMock.
		EXPECT().
		QueryContext(
			ctx,
			"SELECT id, date, number, status FROM sale_order WHERE id = ?",
			saleOrder.ID,
		).
		Return(nil, queryError)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, getErr, queryError.Error())
}

func TestGetByID_ScanError(t *testing.T) {
	// arrange
	ctx := context.Background()

	db, mock, _ := sqlmock.New()

	repository := NewRepository(db)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			ID:     100,
			Number: "123",
			Date:   time.Now().Truncate(time.Second),
			Status: document.StatusDraft,
		},
	}

	rowsResult := sqlmock.NewRows([]string{
		"id",
		"date",
		"number",
	}).AddRow(
		saleOrder.ID,
		helpers.TimeToString(saleOrder.Date),
		saleOrder.Number,
	)

	mock.
		ExpectQuery("SELECT id, date, number, status FROM sale_order WHERE id = ?").
		WithArgs(saleOrder.ID).
		WillReturnRows(rowsResult)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, getErr, "destination arguments in Scan")
}

func TestGetByID_NotFound(t *testing.T) {
	// arrange
	ctx := context.Background()

	db, mock, _ := sqlmock.New()

	repository := NewRepository(db)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			ID:     100,
			Number: "123",
			Date:   time.Now().Truncate(time.Second),
			Status: document.StatusDraft,
		},
	}

	rowsResult := sqlmock.NewRows([]string{
		"id",
		"date",
		"number",
		"status",
	})

	mock.
		ExpectQuery("SELECT id, date, number, status FROM sale_order WHERE id = ?").
		WithArgs(saleOrder.ID).
		WillReturnRows(rowsResult)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.NoError(t, getErr)
}

func TestGetByID_BadDate(t *testing.T) {
	// arrange
	ctx := context.Background()

	db, mock, _ := sqlmock.New()

	repository := NewRepository(db)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			ID:     100,
			Number: "123",
			Date:   time.Now().Truncate(time.Second),
			Status: document.StatusDraft,
		},
	}

	rowsResult := sqlmock.NewRows([]string{
		"id",
		"date",
		"number",
		"status",
	}).AddRow(
		saleOrder.ID,
		"0000-00-00",
		saleOrder.Number,
		saleOrder.Status,
	)

	mock.
		ExpectQuery("SELECT id, date, number, status FROM sale_order WHERE id = ?").
		WithArgs(saleOrder.ID).
		WillReturnRows(rowsResult)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, getErr, "bad date: 0000-00-00")
}

func TestGetByID_BadStatus(t *testing.T) {
	// arrange
	ctx := context.Background()

	db, mock, _ := sqlmock.New()

	repository := NewRepository(db)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			ID:     100,
			Number: "123",
			Date:   time.Now().Truncate(time.Second),
			Status: document.StatusDraft,
		},
	}

	rowsResult := sqlmock.NewRows([]string{
		"id",
		"date",
		"number",
		"status",
	}).AddRow(
		saleOrder.ID,
		helpers.TimeToString(saleOrder.Date),
		saleOrder.Number,
		999,
	)

	mock.
		ExpectQuery("SELECT id, date, number, status FROM sale_order WHERE id = ?").
		WithArgs(saleOrder.ID).
		WillReturnRows(rowsResult)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, getErr, "bad status: 999")
}
