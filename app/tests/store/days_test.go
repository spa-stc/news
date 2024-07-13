package storetests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/spa-stc/newsletter/store"
	"github.com/stretchr/testify/assert"
)

func TestStoreDays(t *testing.T) {
	ctx := context.Background()
	s := NewTestingStore(ctx, t)

	fixture, err := s.UpsertDay(ctx, getDaysFixture())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fixture)

	t.Run("test_days_upsert", func(t *testing.T) {
		fixture.ApInfo = "hi"

		_, err := s.UpsertDay(ctx, fixture)

		assert.NoError(t, err)
	})

	t.Run("test_days_query_success", func(t *testing.T) {
		data, err := s.FindDay(ctx, store.FindDay{
			Date: fixture.Date,
		})

		assert.NoError(t, err)

		assert.Equal(t, fixture, data)
	})

	t.Run("test_days_query_failiure", func(t *testing.T) {
		_, err := s.FindDay(ctx, store.FindDay{
			Date: "hello",
		})

		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("test_week_query_success", func(t *testing.T) {
		dates, days := getDaysFixtures()

		var d []store.Day
		for _, day := range days {
			day, err := s.UpsertDay(ctx, day)
			assert.NoError(t, err)

			d = append(d, day)
		}

		queried, err := s.FindDays(ctx, store.FindDays{
			Dates: dates,
		})

		assert.NoError(t, err)

		assert.Equal(t, queried, d)
	})
}

func getDaysFixture() store.Day {
	return store.Day{
		Date: "2006-12-18",

		Lunch:   "Teriyaki Chicken",
		XPeriod: "Senior Speeches In HUSS",
	}
}

func getDaysFixtures() ([]string, []store.Day) {
	days := []store.Day{
		{
			Date: "2006-01-02",
		},
		{
			Date: "2006-01-03",
		},
		{
			Date: "2006-01-04",
		},
	}

	dates := []string{"2006-01-02", "2006-01-03", "2006-01-04"}

	return dates, days
}
