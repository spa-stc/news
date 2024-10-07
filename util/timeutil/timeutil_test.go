package timeutil_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/util/timeutil"
)

func TestWeek(t *testing.T) {
	t.Parallel()

	start := time.Sunday
	today := time.Date(2023, 12, 3, 23, 0, 0, 0, time.UTC)

	expected := []time.Time{
		time.Date(2023, 12, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 4, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 6, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 7, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 8, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 9, 0, 0, 0, 0, time.UTC),
	}

	actual := timeutil.GetWeek(start, today)

	require.Equal(t, expected, actual)
}

func TestIsWeekday(t *testing.T) {
	t.Parallel()

	w := timeutil.IsWeekday(time.Date(2023, 12, 18, 0, 0, 0, 0, time.UTC))

	require.True(t, w)

	w = timeutil.IsWeekday(time.Date(2023, 12, 22, 0, 0, 0, 0, time.UTC))

	require.True(t, w)

	w = timeutil.IsWeekday(time.Date(2023, 12, 16, 0, 0, 0, 0, time.UTC))

	require.False(t, w)
}
