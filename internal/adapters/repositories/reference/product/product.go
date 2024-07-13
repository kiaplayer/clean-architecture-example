package product

import (
	"context"
	"database/sql"

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

func (r *Repository) Exists(ctx context.Context, id uint64) (bool, error) {
	queryResult, err := r.DB(ctx).QueryContext(
		ctx,
		"SELECT id FROM product WHERE id = ?",
		id,
	)
	if err != nil {
		return false, err
	}

	defer func(queryResult *sql.Rows) {
		_ = queryResult.Close()
	}(queryResult)

	if queryResult.Next() {
		return true, nil
	}

	return false, queryResult.Err()
}
