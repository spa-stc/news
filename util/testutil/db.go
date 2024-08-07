package testutil

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/db"
)

var (
	dbPool   *pgxpool.Pool //nolint:gochecknoglobals // Ok since only used in testing.
	dbPoolMu sync.RWMutex  //nolint:gochecknoglobals // Ok since only used in testing.
)

// Get a test transaction that rolls back after running.
func TestTx(ctx context.Context, t testing.TB) db.Tx {
	tryPool := func() *pgxpool.Pool {
		dbPoolMu.RLock()
		defer dbPoolMu.RUnlock()
		return dbPool
	}

	getPool := func() *pgxpool.Pool {
		if pool := tryPool(); pool != nil {
			return pool
		}

		dbPoolMu.Lock()
		defer dbPoolMu.Unlock()

		var err error
		dbPool, err = pgxpool.New(ctx, os.Getenv("NEWSLETTER_DATABASE_URL"))
		require.NoError(t, err)

		return dbPool
	}

	tx, err := getPool().Begin(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := tx.Rollback(ctx)

		if !errors.Is(err, pgx.ErrTxClosed) {
			require.NoError(t, err)
		}
	})

	return tx
}
