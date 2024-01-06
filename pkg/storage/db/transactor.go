package db

import (
	"context"
	"database/sql"
)

type Transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) *Transactor {
	return &Transactor{
		db: db,
	}
}

type txKey struct{}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func (t *Transactor) RunInTx(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var done bool

	defer func() {
		if !done {
			_ = tx.Rollback()
		}
	}()

	result, err := fn(injectTx(ctx, tx))
	if err != nil {
		return nil, err
	}

	done = true
	return result, tx.Commit()
}
