package sqlite

import (
	"context"
	"embed"

	"github.com/pkg/errors"
)

//go:embed migrations/*.sql
var migrations embed.FS

func (d *DB) Migrate(ctx context.Context) error {
	buf, err := migrations.ReadFile("migrations/SCHEMA.sql")
	if err != nil {
		return errors.Wrap(err, "error reading schema file")
	}

	stmt := string(buf)
	if err := d.execute(ctx, stmt); err != nil {
		return errors.Wrap(err, "error applying schema to database")
	}

	return nil
}

// execute runs a single SQL statement within a transaction.
func (d *DB) execute(ctx context.Context, stmt string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, stmt); err != nil {
		return errors.Wrap(err, "failed to execute statement")
	}

	return tx.Commit()
}
