package sale_order

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/reference"
	"github.com/kiaplayer/clean-architecture-example/pkg/helpers"
	"github.com/kiaplayer/clean-architecture-example/pkg/storage/db"
)

type Repository struct {
	*db.TransactionalRepository
}

func NewRepository(qe db.QueryExecutor) *Repository {
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

	for i, product := range order.Products {
		insertResult, err := r.DB(ctx).ExecContext(
			ctx,
			"INSERT INTO sale_order_product (parent_id, product_id, quantity, price) VALUES (?, ?, ?, ?)",
			order.ID,
			product.Product.ID,
			product.Quantity,
			product.Price,
		)
		if err != nil {
			return nil, err
		}

		lastID, err := insertResult.LastInsertId()
		if err != nil {
			return nil, err
		}

		order.Products[i].ID = uint64(lastID)
	}

	return order, nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*document.SaleOrder, error) {
	saleOrderDTO := struct {
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

	defer func(queryResult *sql.Rows) {
		_ = queryResult.Close()
	}(queryResult)

	if queryResult.Next() {
		err = queryResult.Scan(&saleOrderDTO.ID, &saleOrderDTO.Date, &saleOrderDTO.Number, &saleOrderDTO.Status)
		if err != nil {
			return nil, err
		}
	} else {
		if queryResult.Err() != nil {
			return nil, queryResult.Err()
		}
		return nil, nil
	}

	date, err := helpers.StringToTime(saleOrderDTO.Date)
	if err != nil {
		return nil, fmt.Errorf("bad date: %s", saleOrderDTO.Date)
	}

	status := document.Status(saleOrderDTO.Status)
	if !slices.Contains(document.ValidStatuses, status) {
		return nil, fmt.Errorf("bad status: %d", status)
	}

	result := &document.SaleOrder{
		Document: document.Document{
			ID:     saleOrderDTO.ID,
			Number: saleOrderDTO.Number,
			Date:   date,
			Status: status,
		},
	}

	saleOrderProductDTO := struct {
		ID            uint64
		ProductID     uint64
		ProductName   string
		ProductStatus int
		Quantity      int
		Price         float32
	}{}

	queryResult, err = r.DB(ctx).QueryContext(
		ctx,
		`
			SELECT 
				sop.id, 
				sop.product_id, 
				sop.quantity, 
				sop.price,
				p.name,
				p.status
			FROM sale_order_product AS sop
			LEFT JOIN product AS p ON p.id = sop.product_id
			WHERE sop.parent_id = ?
		`,
		result.ID,
	)
	if err != nil {
		return nil, err
	}

	defer func(queryResult *sql.Rows) {
		_ = queryResult.Close()
	}(queryResult)

	for queryResult.Next() {
		err = queryResult.Scan(
			&saleOrderProductDTO.ID,
			&saleOrderProductDTO.ProductID,
			&saleOrderProductDTO.Quantity,
			&saleOrderProductDTO.Price,
			&saleOrderProductDTO.ProductName,
			&saleOrderProductDTO.ProductStatus,
		)
		if err != nil {
			return nil, err
		}

		productStatus := reference.Status(saleOrderProductDTO.ProductStatus)
		if !slices.Contains(reference.ValidStatuses, productStatus) {
			return nil, fmt.Errorf("bad status: %d", productStatus)
		}

		result.Products = append(result.Products, document.SaleOrderProduct{
			ID: saleOrderProductDTO.ID,
			Product: reference.Product{
				Reference: reference.Reference{
					ID:     saleOrderProductDTO.ProductID,
					Name:   saleOrderProductDTO.ProductName,
					Status: productStatus,
				},
			},
			Quantity: saleOrderProductDTO.Quantity,
			Price:    saleOrderProductDTO.Price,
		})
	}

	if queryResult.Err() != nil {
		return nil, queryResult.Err()
	}

	return result, nil
}
