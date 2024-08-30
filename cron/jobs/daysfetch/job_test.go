package daysfetch_test

import (
	"context"
	"testing"
	"time"

	"github.com/neilotoole/slogt"
	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/cron/jobs/daysfetch"
	"stpaulacademy.tech/newsletter/resource"
	"stpaulacademy.tech/newsletter/util/sliceutil"
	"stpaulacademy.tech/newsletter/util/testutil"
)

type TestGetter struct {
	lunches map[string]string
	info    map[string]daysfetch.CsvData
}

func (t *TestGetter) GetLunch(_ context.Context) (map[string]string, error) {
	return t.lunches, nil
}

func (t *TestGetter) GetInfo(_ context.Context) (map[string]daysfetch.CsvData, error) {
	return t.info, nil
}

func buildFixtures(dates []time.Time) (map[string]string, map[string]daysfetch.CsvData, []resource.Day) {
	lunch := make(map[string]string)
	for _, date := range dates {
		lunch[date.Format(time.DateOnly)] = "lunch"
	}

	other := make(map[string]daysfetch.CsvData)
	for _, date := range dates {
		other[date.Format(daysfetch.CSVDateFormat)] = daysfetch.CsvData{
			Date:     date.Format(daysfetch.CSVDateFormat),
			Rday:     "rotation_day",
			Location: "location",
			Event:    "x_period",
			Grade9:   "grade_9",
			Grade10:  "grade_10",
			Grade11:  "grade_11",
			Grade12:  "grade_12",
			ApInfo:   "ap_info",
			CcInfo:   "cc_info",
		}
	}

	days := sliceutil.Map(dates, func(d time.Time) resource.Day {
		return resource.Day{
			Date:        d,
			Lunch:       "lunch",
			XPeriod:     "x_period",
			RotationDay: "rotation_day",
			Location:    "location",
			ApInfo:      "ap_info",
			CcInfo:      "cc_info",
			Grade9:      "grade_9",
			Grade10:     "grade_10",
			Grade11:     "grade_11",
			Grade12:     "grade_12",
			CreatedTS:   time.UnixMicro(0),
			UpdatedTS:   time.UnixMicro(0),
		}
	})

	return lunch, other, days
}

func TestDaysUpdateJob(t *testing.T) {
	t.Parallel()
	logger := slogt.New(t)

	dates := []time.Time{
		time.Date(2023, 12, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 4, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 6, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 7, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 8, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 9, 0, 0, 0, 0, time.UTC),
	}

	ctx := context.Background()

	executor := testutil.TestTx(ctx, t)

	timegen := testutil.NewTestTimeGen(time.Date(2023, 12, 3, 23, 0, 0, 0, time.UTC))

	lunches, csvdata, expected := buildFixtures(dates)

	getter := &TestGetter{
		info:    csvdata,
		lunches: lunches,
	}

	job := daysfetch.New(executor, timegen, getter, nil, false)

	err := job.Run(ctx)
	require.NoError(t, err)

	week, err := resource.GetManyDays(ctx, executor, dates)
	require.NoError(t, err)

	week = sliceutil.Map(week, func(r resource.Day) resource.Day {
		r.CreatedTS = time.UnixMicro(0)
		r.UpdatedTS = time.UnixMicro(0)

		return r
	})

	logger.InfoContext(ctx, "fetched", "data", week)
	logger.InfoContext(ctx, "expected", "data", expected)

	require.ElementsMatch(t, expected, week)
}
