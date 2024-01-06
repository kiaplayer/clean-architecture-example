package db

import (
	"context"
)

type TransactionalRepository struct {
	db QueryExecutor
}

func NewTransactionalRepository(db QueryExecutor) *TransactionalRepository {
	return &TransactionalRepository{
		db: db,
	}
}

func (r *TransactionalRepository) DB(ctx context.Context) QueryExecutor {
	tx := extractTx(ctx)
	if tx != nil {
		return tx
	}
	return r.db
}
