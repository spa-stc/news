package storetests

import (
	"context"
	"testing"

	"github.com/spa-stc/newsletter/store"
	"github.com/spa-stc/newsletter/store/db"
	"github.com/spa-stc/newsletter/tests"
)

func NewTestingStore(ctx context.Context, t *testing.T) *store.Store {
	profile := tests.GetTestingProfile(t)
	driver, err := db.GetDriver(profile)
	if err != nil {
		t.Fatal(err)
	}

	if err := driver.Migrate(ctx); err != nil {
		t.Fatal(err)
	}

	return store.New(profile, driver)
}
