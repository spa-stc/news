package store

import (
	"context"
	"time"

	"github.com/spa-stc/newsletter/timeutil"
)

const DayFormat = "2006-01-02"

// Main Day Model, Outputted From Cron Jobs.
type Day struct {
	// Day format, 2006-01-02.
	Date,

	Lunch,
	XPeriod,
	RotationDay,
	Location,
	Notes,
	ApInfo,
	CCInfo,
	Grade9,
	Grade10,
	Grade11,
	Grade12 string

	CreatedTs int64
	UpdatedTs int64
}

type FindDay struct {
	Date string
}

type FindDays struct {
	Dates []string
}

func (s *Store) UpsertDay(ctx context.Context, day Day) (Day, error) {
	return s.db.UpsertDay(ctx, day)
}

func (s *Store) FindDay(ctx context.Context, query FindDay) (Day, error) {
	return s.db.FindDay(ctx, query)
}

func (s *Store) FindDays(ctx context.Context, query FindDays) ([]Day, error) {
	return s.db.FindDays(ctx, query)
}

func (s *Store) GetWeek(ctx context.Context) ([]Day, error) {
	week := timeutil.GetWeek(time.Now())
	var weekKeys []string
	for _, day := range week {
		weekKeys = append(weekKeys, day.Format(DayFormat))
	}

	return s.FindDays(ctx, FindDays{
		Dates: weekKeys,
	})
}
