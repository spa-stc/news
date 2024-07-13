package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/pkg/errors"
	"github.com/spa-stc/newsletter/server/profile"
	"github.com/spa-stc/newsletter/store"
)

type DB struct {
	db      *sql.DB
	profile *profile.Profile
}

func NewDB(profile *profile.Profile) (store.Driver, error) {
	if profile.DSN == "" {
		return nil, errors.New("missing database name")
	}

	dbName := fmt.Sprintf("%s/%s", profile.Dir, profile.DSN)

	sqliteDB, err := sql.Open("sqlite", dbName+"?_pragma=foreign_keys(0)&_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to database")
	}

	driver := DB{db: sqliteDB, profile: profile}

	return &driver, nil
}

func (d *DB) GetDB() *sql.DB {
	return d.db
}

func (d *DB) Close() error {
	return d.db.Close()
}

func getPlaceholders(len uint) []string {
	return make([]string, len, '?')
}
