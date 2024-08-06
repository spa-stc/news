package resource

import (
	"context"
	"time"

	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/db/dbsqlc"
	"stpaulacademy.tech/newsletter/util/sliceutil"
)

type Day struct {
	Date        time.Time `json:"date"`
	Lunch       string    `json:"lunch"`
	XPeriod     string    `json:"x_period"`
	RotationDay string    `json:"r_day"`
	Location    string    `json:"location"`
	Notes       string    `json:"notes"`
	ApInfo      string    `json:"ap_info"`
	CcInfo      string    `json:"cc_info"`
	Grade9      string    `json:"grade_9"`
	Grade10     string    `json:"grade_10"`
	Grade11     string    `json:"grade_11"`
	Grade12     string    `json:"grade_12"`
	CreatedTS   time.Time `json:"created_ts"`
	UpdatedTS   time.Time `json:"updated_ts"`
}

func fromSqlcDay(d dbsqlc.Day) Day {
	return Day{
		Date:        d.Date.Local(),
		Lunch:       d.Lunch,
		XPeriod:     d.XPeriod,
		RotationDay: d.RotationDay,
		Location:    d.Location,
		Notes:       d.Notes,
		ApInfo:      d.ApInfo,
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

	return sliceutil.Map(days, fromSqlcDay), nil
}

func BatchUpsertDays(ctx context.Context, e db.Executor, days []Day) error {
	sqlc := dbsqlc.New()

	err := db.ExecInTx(ctx, e, func(tx db.Tx) error {
		for _, d := range days {
			if _, err := sqlc.UpsertDay(ctx, tx, dbsqlc.UpsertDayParams{
				Date:        d.Date.Local(),
				Lunch:       d.Lunch,
				XPeriod:     d.XPeriod,
				RotationDay: d.RotationDay,
				Location:    d.Location,
				Notes:       d.Notes,
				ApInfo:      d.ApInfo,
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
