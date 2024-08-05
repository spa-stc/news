package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Generic database executor that can be passed around the application.
type Executor interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row

	Begin(ctx context.Context) (pgx.Tx, error)
}

// Specific, transaction based executor that allows for committing and rolling back transactions.
type Tx interface {
	Executor

	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func ExecInTx(ctx context.Context, executor Executor, fun func(Tx) error) error {
	newfun := func(t Tx) (struct{}, error) {
		err := fun(t)
		return struct{}{}, err
	}

	_, err := QueryInTx(ctx, executor, newfun)
	return err
}

func QueryInTx[R any](ctx context.Context, executor Executor, fun func(Tx) (R, error)) (R, error) {
	tx, err := executor.Begin(ctx)
	if err != nil {
		var result R
		return result, fmt.Errorf("error beginning transaction: %w", err)
	}
	defer func() {
		if e := tx.Rollback(ctx); e != nil && !errors.Is(e, pgx.ErrTxClosed) {
			slog.Error("error rolling back transaction", "error", e)
		}
	}()

	value, err := fun(tx)
	if err != nil {
		var result R
		return result, err
	}

	if err := tx.Commit(ctx); err != nil {
		var result R
		return result, fmt.Errorf("error committing transaction: %w", err)
	}

	return value, nil
}
