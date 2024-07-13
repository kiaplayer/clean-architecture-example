package sale_order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/reference"
	"github.com/kiaplayer/clean-architecture-example/pkg/helpers"
)

func TestCreateOrder_Success(t *testing.T) {
	// arrange
	ctx := context.Background()

	db, mock, _ := sqlmock.New()

	repository := NewRepository(db)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Number: "123",
			Date:   time.Now(),
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

	saleOrderID := int64(100)
	saleOrderProductID := int64(200)

	insertSaleOrderResult := sqlmock.NewResult(saleOrderID, 1)
	insertSaleOrderProductResult := sqlmock.NewResult(saleOrderProductID, 1)

	mock.
		ExpectExec("INSERT INTO sale_order").
		WithArgs(
			saleOrder.Number,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Status,
		).
		WillReturnResult(insertSaleOrderResult)

	mock.
		ExpectExec("INSERT INTO sale_order_product").
		WithArgs(
			saleOrderID,
			saleOrder.Products[0].Product.ID,
			saleOrder.Products[0].Quantity,
			saleOrder.Products[0].Price,
		).
		WillReturnResult(insertSaleOrderProductResult)

	// act
	updatedSaleOrder, createErr := repository.CreateOrder(ctx, saleOrder)

	// assert
	assert.NoError(t, createErr)
	assert.NotNil(t, updatedSaleOrder)
	assert.Equal(t, updatedSaleOrder.ID, uint64(saleOrderID))
	assert.Equal(t, updatedSaleOrder.Products[0].ID, uint64(saleOrderProductID))
}

func TestCreateOrder_InsertError(t *testing.T) {
	// arrange
	ctx := context.Background()

	db, mock, _ := sqlmock.New()

	repository := NewRepository(db)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Number: "123",
			Date:   time.Now(),
			Status: document.StatusDraft,
		},
	}

	insertError := errors.New("insert error")

	mock.
		ExpectExec("INSERT INTO sale_order").
		WithArgs(
			saleOrder.Number,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Status,
		).
		WillReturnError(insertError)

	// act
	updatedSaleOrder, createErr := repository.CreateOrder(ctx, saleOrder)

	// assert
	assert.ErrorContains(t, createErr, insertError.Error())
	assert.Nil(t, updatedSaleOrder)
}

func TestCreateOrder_InsertProductError(t *testing.T) {
	// arrange
	ctx := context.Background()

	db, mock, _ := sqlmock.New()

	repository := NewRepository(db)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Number: "123",
			Date:   time.Now(),
			Status: document.StatusDraft,
		},
		Products: []document.SaleOrderProduct{
			{
				ID: 1000,
				Product: reference.Product{
					Reference: reference.Reference{
						ID:     1,
						Name:   "Keyboard",
						Status: reference.StatusActive,
					},
				},
				Quantity: 1,
				Price:    300,
			},
		},
	}

	saleOrderID := int64(100)

	insertSaleOrderResult := sqlmock.NewResult(saleOrderID, 1)
	insertError := errors.New("insert product error")

	mock.
		ExpectExec("INSERT INTO sale_order").
		WithArgs(
			saleOrder.Number,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Status,
		).
		WillReturnResult(insertSaleOrderResult)

	mock.
		ExpectExec("INSERT INTO sale_order_product").
		WithArgs(
			saleOrderID,
			saleOrder.Products[0].Product.ID,
			saleOrder.Products[0].Quantity,
			saleOrder.Products[0].Price,
		).
		WillReturnError(insertError)

	// act
	updatedSaleOrder, createErr := repository.CreateOrder(ctx, saleOrder)

	// assert
	assert.ErrorContains(t, createErr, insertError.Error())
	assert.Nil(t, updatedSaleOrder)
}

func TestCreateOrder_LastInsertIDError(t *testing.T) {
	// arrange
	ctx := context.Background()

	db, mock, _ := sqlmock.New()

	repository := NewRepository(db)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Number: "123",
			Date:   time.Now(),
			Status: document.StatusDraft,
		},
	}

	lastInsertIDError := errors.New("last insert id error")

	insertResult := sqlmock.NewErrorResult(lastInsertIDError)

	mock.
		ExpectExec("INSERT INTO sale_order").
		WithArgs(
			saleOrder.Number,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Status,
		).
		WillReturnResult(insertResult)

	// act
	updatedSaleOrder, createErr := repository.CreateOrder(ctx, saleOrder)

	// assert
	assert.ErrorContains(t, createErr, lastInsertIDError.Error())
	assert.Nil(t, updatedSaleOrder)
}

func TestCreateOrder_LastInsertProductIDError(t *testing.T) {
	// arrange
	ctx := context.Background()

	db, mock, _ := sqlmock.New()

	repository := NewRepository(db)

	saleOrder := &document.SaleOrder{
		Document: document.Document{
			Number: "123",
			Date:   time.Now(),
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

	saleOrderID := int64(100)
	insertSaleOrderResult := sqlmock.NewResult(saleOrderID, 1)

	lastInsertProductIDError := errors.New("last insert id error")
	insertSaleOrderProductResult := sqlmock.NewErrorResult(lastInsertProductIDError)

	mock.
		ExpectExec("INSERT INTO sale_order").
		WithArgs(
			saleOrder.Number,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Status,
		).
		WillReturnResult(insertSaleOrderResult)

	mock.
		ExpectExec("INSERT INTO sale_order_product").
		WithArgs(
			saleOrderID,
			saleOrder.Products[0].Product.ID,
			saleOrder.Products[0].Quantity,
			saleOrder.Products[0].Price,
		).
		WillReturnResult(insertSaleOrderProductResult)

	// act
	updatedSaleOrder, createErr := repository.CreateOrder(ctx, saleOrder)

	// assert
	assert.ErrorContains(t, createErr, lastInsertProductIDError.Error())
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
		Products: []document.SaleOrderProduct{
			{
				ID: 1000,
				Product: reference.Product{
					Reference: reference.Reference{
						ID:     1,
						Name:   "Keyboard",
						Status: reference.StatusActive,
					},
				},
				Quantity: 1,
				Price:    300,
			},
		},
	}

	saleOrdersResult := sqlmock.NewRows([]string{
		"id",
		"date",
		"number",
		"status",
	}).
		AddRow(
			saleOrder.ID,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Number,
			saleOrder.Status,
		)

	saleOrdersProductsResult := sqlmock.NewRows([]string{
		"id",
		"product_id",
		"quantity",
		"price",
		"name",
		"status",
	}).
		AddRow(
			saleOrder.Products[0].ID,
			saleOrder.Products[0].Product.ID,
			saleOrder.Products[0].Quantity,
			saleOrder.Products[0].Price,
			saleOrder.Products[0].Product.Name,
			saleOrder.Products[0].Product.Status,
		)

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order ").
		WithArgs(saleOrder.ID).
		WillReturnRows(saleOrdersResult)

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order_product ").
		WithArgs(saleOrder.ID).
		WillReturnRows(saleOrdersProductsResult)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.NoError(t, getErr)
	assert.Equal(t, saleOrder, actualSaleOrder)
}

func TestGetByID_QueryError(t *testing.T) {
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

	queryError := errors.New("query error")

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order ").
		WithArgs(saleOrder.ID).
		WillReturnError(queryError)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, getErr, queryError.Error())
}

func TestGetByID_QueryProductsError(t *testing.T) {
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
		Products: []document.SaleOrderProduct{
			{
				ID: 1000,
				Product: reference.Product{
					Reference: reference.Reference{
						ID:     1,
						Name:   "Keyboard",
						Status: reference.StatusActive,
					},
				},
				Quantity: 1,
				Price:    300,
			},
		},
	}

	saleOrdersResult := sqlmock.NewRows([]string{
		"id",
		"date",
		"number",
		"status",
	}).
		AddRow(
			saleOrder.ID,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Number,
			saleOrder.Status,
		)

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order ").
		WithArgs(saleOrder.ID).
		WillReturnRows(saleOrdersResult)

	queryError := errors.New("query error")

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order_product ").
		WithArgs(saleOrder.ID).
		WillReturnError(queryError)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, getErr, queryError.Error())
}

func TestGetByID_NextError(t *testing.T) {
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

	nextError := errors.New("db next error")

	rowsResult := sqlmock.NewRows([]string{
		"id",
		"date",
		"number",
		"status",
	}).
		AddRow(
			saleOrder.ID,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Number,
			saleOrder.Status,
		).
		RowError(0, nextError)

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order ").
		WithArgs(saleOrder.ID).
		WillReturnRows(rowsResult)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, getErr, nextError.Error())
}

func TestGetByID_NextProductError(t *testing.T) {
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
		Products: []document.SaleOrderProduct{
			{
				ID: 1000,
				Product: reference.Product{
					Reference: reference.Reference{
						ID:     1,
						Name:   "Keyboard",
						Status: reference.StatusActive,
					},
				},
				Quantity: 1,
				Price:    300,
			},
		},
	}

	nextError := errors.New("db next error")

	saleOrderResult := sqlmock.NewRows([]string{
		"id",
		"date",
		"number",
		"status",
	}).
		AddRow(
			saleOrder.ID,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Number,
			saleOrder.Status,
		)

	saleOrdersProductsResult := sqlmock.NewRows([]string{
		"id",
		"product_id",
		"quantity",
		"price",
		"name",
		"status",
	}).
		AddRow(
			saleOrder.Products[0].ID,
			saleOrder.Products[0].Product.ID,
			saleOrder.Products[0].Quantity,
			saleOrder.Products[0].Price,
			saleOrder.Products[0].Product.Name,
			saleOrder.Products[0].Product.Status,
		).
		RowError(0, nextError)

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order ").
		WithArgs(saleOrder.ID).
		WillReturnRows(saleOrderResult)

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order_product ").
		WithArgs(saleOrder.ID).
		WillReturnRows(saleOrdersProductsResult)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, getErr, nextError.Error())
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
	}).
		AddRow(
			saleOrder.ID,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Number,
		)

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order ").
		WithArgs(saleOrder.ID).
		WillReturnRows(rowsResult)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, getErr, "destination arguments in Scan")
}

func TestGetByID_ScanProductsError(t *testing.T) {
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
		Products: []document.SaleOrderProduct{
			{
				ID: 1000,
				Product: reference.Product{
					Reference: reference.Reference{
						ID:     1,
						Name:   "Keyboard",
						Status: reference.StatusActive,
					},
				},
				Quantity: 1,
				Price:    300,
			},
		},
	}

	saleOrdersResult := sqlmock.NewRows([]string{
		"id",
		"date",
		"number",
		"status",
	}).
		AddRow(
			saleOrder.ID,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Number,
			saleOrder.Status,
		)

	saleOrdersProductsResult := sqlmock.NewRows([]string{
		"id",
		"product_id",
		"quantity",
		"price",
		"name",
	}).
		AddRow(
			saleOrder.Products[0].ID,
			saleOrder.Products[0].Product.ID,
			saleOrder.Products[0].Quantity,
			saleOrder.Products[0].Price,
			saleOrder.Products[0].Product.Name,
		)

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order ").
		WithArgs(saleOrder.ID).
		WillReturnRows(saleOrdersResult)

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order_product ").
		WithArgs(saleOrder.ID).
		WillReturnRows(saleOrdersProductsResult)

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
		ExpectQuery("^SELECT (.+) FROM sale_order ").
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
	}).
		AddRow(
			saleOrder.ID,
			"0000-00-00",
			saleOrder.Number,
			saleOrder.Status,
		)

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order ").
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
	}).
		AddRow(
			saleOrder.ID,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Number,
			999,
		)

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order ").
		WithArgs(saleOrder.ID).
		WillReturnRows(rowsResult)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, getErr, "bad status: 999")
}

func TestGetByID_BadProductStatus(t *testing.T) {
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
		Products: []document.SaleOrderProduct{
			{
				ID: 1000,
				Product: reference.Product{
					Reference: reference.Reference{
						ID:     1,
						Name:   "Keyboard",
						Status: reference.StatusActive,
					},
				},
				Quantity: 1,
				Price:    300,
			},
		},
	}

	saleOrderResult := sqlmock.NewRows([]string{
		"id",
		"date",
		"number",
		"status",
	}).
		AddRow(
			saleOrder.ID,
			helpers.TimeToString(saleOrder.Date),
			saleOrder.Number,
			saleOrder.Status,
		)

	saleOrdersProductsResult := sqlmock.NewRows([]string{
		"id",
		"product_id",
		"quantity",
		"price",
		"name",
		"status",
	}).
		AddRow(
			saleOrder.Products[0].ID,
			saleOrder.Products[0].Product.ID,
			saleOrder.Products[0].Quantity,
			saleOrder.Products[0].Price,
			saleOrder.Products[0].Product.Name,
			999,
		)

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order ").
		WithArgs(saleOrder.ID).
		WillReturnRows(saleOrderResult)

	mock.
		ExpectQuery("^SELECT (.+) FROM sale_order_product ").
		WithArgs(saleOrder.ID).
		WillReturnRows(saleOrdersProductsResult)

	// act
	actualSaleOrder, getErr := repository.GetByID(ctx, saleOrder.ID)

	// assert
	assert.Nil(t, actualSaleOrder)
	assert.ErrorContains(t, getErr, "bad status: 999")
}
