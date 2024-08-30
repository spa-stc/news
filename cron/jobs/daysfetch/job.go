package daysfetch

import (
	"context"
	"fmt"
	"time"

	"stpaulacademy.tech/newsletter/cron"
	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/resource"
	"stpaulacademy.tech/newsletter/util/service"
	"stpaulacademy.tech/newsletter/util/sliceutil"
	"stpaulacademy.tech/newsletter/util/timeutil"
)

type Job struct {
	db       db.Executor
	timeGen  service.TimeGenerator
	getter   DataGetter
	notifer  cron.StatusNotifer
	schedule string
}

func New(
	db db.Executor,
	timeGen service.TimeGenerator,
	getter DataGetter,
	notifer cron.StatusNotifer,
	development bool,
) *Job {
	schedule := "5 0 * * *"
	if development {
		schedule = "* * * * *"
	}

	return &Job{
		db:       db,
		timeGen:  timeGen,
		getter:   getter,
		notifer:  notifer,
		schedule: schedule,
	}
}

func (j *Job) Run(ctx context.Context) error {
	dates := timeutil.GetWeek(time.Sunday, j.timeGen.NowUTC())

	lunches, err := j.getter.GetLunch(ctx)
	if err != nil {
		return fmt.Errorf("error getting lunch: %w", err)
	}

	otherInfo, err := j.getter.GetInfo(ctx)
	if err != nil {
		return fmt.Errorf("error getting other info: %w", err)
	}

	days := sliceutil.Map(dates, func(t time.Time) resource.Day {
		day := resource.Day{
			Date: t,
		}

		if lunch, ok := lunches[t.Format(time.DateOnly)]; ok {
			day.Lunch = lunch
		} else {
			day.Lunch = "Not Available"
		}

		if info, ok := otherInfo[t.Format(CSVDateFormat)]; ok {
			day.RotationDay = info.Rday
			day.Location = info.Location
			day.XPeriod = info.Event
			day.Location = info.Location
			day.Grade9 = info.Grade9
			day.Grade10 = info.Grade10
			day.Grade11 = info.Grade11
			day.Grade12 = info.Grade12
			day.ApInfo = info.ApInfo
			day.CcInfo = info.CcInfo
		}

		return day
	})

	err = resource.BatchUpsertDays(ctx, j.db, days)
	if err != nil {
		return fmt.Errorf("error inserting days: %w", err)
	}

	return nil
}

func (j *Job) Notifer() cron.StatusNotifer {
	return j.notifer
}

func (j *Job) Spec() string {
	return j.schedule
}
