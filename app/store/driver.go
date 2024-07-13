package store

import (
	"context"
	"database/sql"
)

// The driver interface defines possible database operations.
type Driver interface {
	GetDB() *sql.DB
	Close() error

	Migrate(ctx context.Context) error

	UpsertDay(context.Context, Day) (Day, error)
	FindDay(context.Context, FindDay) (Day, error)
	FindDays(context.Context, FindDays) ([]Day, error)
}
