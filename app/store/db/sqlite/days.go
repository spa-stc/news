package sqlite

import (
	"context"

	"github.com/spa-stc/newsletter/store"
)

func (db *DB) UpsertDay(ctx context.Context, day store.Day) (store.Day, error) {
	stmt := `
		INSERT INTO days (
			date,
			lunch, 
			x_period, 
			rotation_day, 
			location, 
			notes,
			ap_info,
			cc_info, 
			grade_9, 
			grade_10, 
			grade_11,
			grade_12
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 
		ON CONFLICT DO UPDATE 
		SET
			date = ?, 
			lunch = ?,
			x_period = ?,
			rotation_day = ?,
			location = ?,
			notes = ?,
			ap_info = ?,
			cc_info = ?,
			grade_9 = ?,
			grade_10 = ?,
			grade_11 = ?,
			grade_12 = ?,
			updated_ts = strftime('%s', 'now') 
		RETURNING 
			created_ts, 
			updated_ts
		`

	args := []any{
		day.Date,
		day.Lunch,
		day.XPeriod,
		day.RotationDay,
		day.Location,
		day.Notes,
		day.ApInfo,
		day.CCInfo,
		day.Grade9,
		day.Grade10,
		day.Grade11,
		day.Grade12,
	}

	args = append(args, args...)

	err := db.db.QueryRowContext(ctx, stmt, args...).
		Scan(&day.CreatedTs, &day.UpdatedTs)
	if err != nil {
		return store.Day{}, err
	}

	return day, nil
}

func (db *DB) FindDay(ctx context.Context, query store.FindDay) (store.Day, error) {
	stmt := `
		SELECT
			date,
			lunch, 
			x_period, 
			rotation_day, 
			location, 
			notes,
			ap_info,
			cc_info, 
			grade_9, 
			grade_10, 
			grade_11,
			grade_12,
			updated_ts,
			created_ts
		FROM 
			days 
		WHERE (
			date = ?
		)
	`

	row := db.db.QueryRowContext(ctx, stmt, query.Date)
	var day store.Day
	if err := row.Scan(
		&day.Date,
		&day.Lunch,
		&day.XPeriod,
		&day.RotationDay,
		&day.Location,
		&day.Notes,
		&day.ApInfo,
		&day.CCInfo,
		&day.Grade9,
		&day.Grade10,
		&day.Grade11,
		&day.Grade12,
		&day.UpdatedTs,
		&day.CreatedTs,
	); err != nil {
		return store.Day{}, err
	}

	return day, nil
}

func (db *DB) FindDays(ctx context.Context, query store.FindDays) ([]store.Day, error) {
	var days []store.Day
	for _, date := range query.Dates {
		day, err := db.FindDay(ctx, store.FindDay{
			Date: date,
		})
		if err != nil {
			return nil, err
		}

		days = append(days, day)
	}

	return days, nil
}
