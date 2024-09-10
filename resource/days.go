package resource

import (
	"context"
	"time"

	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/db/dbsqlc"
	"stpaulacademy.tech/newsletter/util/sliceutil"
)

type Day struct {
	Date        time.Time
	Lunch       string
	XPeriod     string
	RotationDay string
	Location    string
	Notes       string
	CcInfo      string
	Grade9      string
	Grade10     string
	Grade11     string
	Grade12     string
	CreatedTS   time.Time
	UpdatedTS   time.Time
}

func fromSqlcDay(d dbsqlc.Day) Day {
	return Day{
		Date:        d.Date.UTC(),
		Lunch:       d.Lunch,
		XPeriod:     d.XPeriod,
		RotationDay: d.RotationDay,
		Location:    d.Location,
		Notes:       d.Notes,
		CcInfo:      d.CcInfo,
		Grade9:      d.Grade9,
		Grade10:     d.Grade10,
		Grade11:     d.Grade11,
		Grade12:     d.Grade12,
		CreatedTS:   d.CreatedTs.UTC(),
		UpdatedTS:   d.UpdatedTs.UTC(),
	}
}

func GetManyDays(ctx context.Context, e db.Executor, dates []time.Time) ([]Day, error) {
	sqlc := dbsqlc.New()

	days, err := sqlc.GetManyDaysByDate(ctx, e, dates)
	if err != nil {
		return nil, db.HandleError(err)
	}

	if len(days) == 0 {
		return nil, db.ErrNotFound
	}

	return sliceutil.Map(days, fromSqlcDay), nil
}

func GetManyDaysByRange(ctx context.Context, e db.Executor, start, end time.Time) ([]Day, error) {
	days, err := dbsqlc.New().GetManyDaysbyDateRange(ctx, e, dbsqlc.GetManyDaysbyDateRangeParams{
		Date:   start,
		Date_2: end,
	})
	if err != nil {
		return nil, db.HandleError(err)
	}
	if len(days) == 0 {
		return nil, db.ErrNotFound
	}

	return sliceutil.Map(days, fromSqlcDay), nil
}

func BatchUpsertDays(ctx context.Context, e db.Executor, days []Day) error {
	sqlc := dbsqlc.New()

	err := db.ExecInTx(ctx, e, func(tx db.Tx) error {
		for _, d := range days {
			if _, err := sqlc.UpsertDay(ctx, tx, dbsqlc.UpsertDayParams{
				Date:        d.Date.UTC(),
				Lunch:       d.Lunch,
				XPeriod:     d.XPeriod,
				RotationDay: d.RotationDay,
				Location:    d.Location,
				Notes:       d.Notes,
				CcInfo:      d.CcInfo,
				Grade9:      d.Grade9,
				Grade10:     d.Grade10,
				Grade11:     d.Grade11,
				Grade12:     d.Grade12,
			}); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return db.HandleError(err)
	}

	return nil
}
