package repository

import (
	"context"
	"database/sql"
)

type RepositoryImpl struct {
	queries *Queries
}

func NewRepository(db DBTX) *RepositoryImpl {
	return &RepositoryImpl{
		queries: New(db),
	}
}

// Transaction methods
func (r *RepositoryImpl) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return r.queries.BeginTx(ctx, opts)
}

func (r *RepositoryImpl) WithTx(tx *sql.Tx) any {
	return &RepositoryImpl{
		queries: r.queries.WithTx(tx),
	}
}
