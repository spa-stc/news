package resource_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/resource"
	"stpaulacademy.tech/newsletter/util/sliceutil"
	"stpaulacademy.tech/newsletter/util/testutil"
)

func parseDate(t *testing.T, s string) time.Time {
	date, err := time.Parse(time.DateOnly, s)
	require.NoError(t, err)
	return date
}

func TestDaysResource(t *testing.T) {
	t.Parallel()
	ctx := testutil.Setup(t)
	tx := testutil.TestTx(ctx, t)

	t.Run("test_query_not_found", func(t *testing.T) {
		_, err := resource.GetManyDays(ctx, tx, []time.Time{
			parseDate(t, "2022-07-03"),
		})
		require.Error(t, err)
	})

	t.Run("test_query_incomplete", func(t *testing.T) {
		_, err := resource.GetManyDays(ctx, tx, []time.Time{
			parseDate(t, "2022-07-03"),
			parseDate(t, "2024-12-18"),
		})
		require.NoError(t, err)
	})

	t.Run("test_query", func(t *testing.T) {
		daybase := resource.Day{
			Lunch:       "lunch",
			XPeriod:     "x_period",
			RotationDay: "rotation_day",
			Location:    "location",
			Notes:       "notes",
			ApInfo:      "ap_info",
			CcInfo:      "cc_info",
			Grade9:      "grade_9",
			Grade10:     "grade_10",
			Grade11:     "grade_11",
			Grade12:     "grade_12",
		}

		dates := []time.Time{
			parseDate(t, "2024-12-18"),
			parseDate(t, "2024-12-19"),
			parseDate(t, "2024-12-20"),
		}

		expected := sliceutil.Map(dates, func(s time.Time) resource.Day {
			daybase.Date = s
			daybase.CreatedTS = s
			daybase.UpdatedTS = s

			return daybase
		})

		days, err := resource.GetManyDays(ctx, tx, dates)

		days = sliceutil.Map(days, func(d resource.Day) resource.Day {
			d.CreatedTS = d.Date
			d.UpdatedTS = d.Date

			return d
		})

		require.NoError(t, err)
		require.Equal(t, expected, days)
	})

	t.Run("test_insert_upsert", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)

		day := resource.Day{
			Date:        parseDate(t, "2023-12-02"),
			Lunch:       "lunch",
			XPeriod:     "x_period",
			RotationDay: "rotation_day",
			Location:    "location",
			Notes:       "notes",
			ApInfo:      "ap_info",
			CcInfo:      "cc_info",
			Grade9:      "grade_9",
			Grade10:     "grade_10",
			Grade11:     "grade_11",
			Grade12:     "grade_12",
			CreatedTS:   time.UnixMicro(0),
			UpdatedTS:   time.UnixMicro(0),
		}

		day2 := resource.Day{
			Date:        parseDate(t, "2023-12-03"),
			Lunch:       "lunch",
			XPeriod:     "x_period",
			RotationDay: "rotation_day",
			Location:    "location",
			Notes:       "notes",
			ApInfo:      "ap_info",
			CcInfo:      "cc_info",
			Grade9:      "grade_9",
			Grade10:     "grade_10",
			Grade11:     "grade_11",
			Grade12:     "grade_12",
			CreatedTS:   time.UnixMicro(0),
			UpdatedTS:   time.UnixMicro(0),
		}

		err := resource.BatchUpsertDays(ctx, tx, []resource.Day{
			day,
			day2,
		})
		require.NoError(t, err)

		day.Lunch = "yo"
		day2.Lunch = "yo"

		err = resource.BatchUpsertDays(ctx, tx, []resource.Day{
			day,
			day2,
		})
		require.NoError(t, err)

		out, err := resource.GetManyDays(ctx, tx, []time.Time{
			day.Date,
			day2.Date,
		})
		require.NoError(t, err)

		out = sliceutil.Map(out, func(r resource.Day) resource.Day {
			r.CreatedTS = time.UnixMicro(0)
			r.UpdatedTS = time.UnixMicro(0)

			return r
		})

		require.ElementsMatch(t, []resource.Day{day, day2}, out)
	})
}
