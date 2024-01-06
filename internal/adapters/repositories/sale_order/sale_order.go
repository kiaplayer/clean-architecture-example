//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=sale_order_mocks.go
package sale_order

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	"github.com/kiaplayer/clean-architecture-example/pkg/helpers"
	"github.com/kiaplayer/clean-architecture-example/pkg/storage/db"
)

type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type Repository struct {
	*db.TransactionalRepository
}

func NewRepository(qe QueryExecutor) *Repository {
	return &Repository{
		TransactionalRepository: db.NewTransactionalRepository(qe),
	}
}

func (r *Repository) CreateOrder(ctx context.Context, order *document.SaleOrder) (*document.SaleOrder, error) {
	insertResult, err := r.DB(ctx).ExecContext(
		ctx,
		"INSERT INTO sale_order (number, date, status) VALUES (?, ?, ?)",
		order.Number,
		helpers.TimeToString(order.Date),
		order.Status,
	)
	if err != nil {
		return nil, err
	}

	lastID, err := insertResult.LastInsertId()
	if err != nil {
		return nil, err
	}
	order.ID = uint64(lastID)

	return order, nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*document.SaleOrder, error) {
	values := struct {
		ID     uint64
		Number string
		Date   string
		Status int
	}{}

	queryResult, err := r.DB(ctx).QueryContext(
		ctx,
		"SELECT id, date, number, status FROM sale_order WHERE id = ?",
		id,
	)
	if err != nil {
		return nil, err
	}

	if queryResult.Next() {
		err = queryResult.Scan(&values.ID, &values.Date, &values.Number, &values.Status)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	date, err := helpers.StringToTime(values.Date)
	if err != nil {
		return nil, fmt.Errorf("bad date: %s", values.Date)
	}

	status := document.Status(values.Status)
	if !slices.Contains(document.ValidStatuses, status) {
		return nil, fmt.Errorf("bad status: %d", status)
	}

	result := &document.SaleOrder{
		Document: document.Document{
			ID:     values.ID,
			Number: values.Number,
			Date:   date,
			Status: status,
		},
	}

	return result, nil
}
