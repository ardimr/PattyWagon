package repository

import (
	"context"
	"database/sql"
	"errors"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type DBTXWithTx interface {
	DBTX
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

func (q *Queries) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	// Check if the underlying db can begin a transaction
	if beginner, ok := q.db.(interface {
		BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	}); ok {
		return beginner.BeginTx(ctx, opts)
	}
	// If q.db is already a transaction, we can't begin another one
	return nil, errors.New("cannot begin transaction: already in a transaction")
}

func (q *Queries) CanBeginTx() bool {
	_, ok := q.db.(interface {
		BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	})
	return ok
}
