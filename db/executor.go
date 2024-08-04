package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Generic database executor that can be passed around the application.
type Executor interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row

	Begin(ctx context.Context) (Tx, error)
}

// Specific, transaction based executor that allows for committing and rolling back transactions.
type Tx interface {
	Executor

	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
