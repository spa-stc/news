package timeutil_test

import (
	"testing"
	"time"

	"github.com/spa-stc/newsletter/timeutil"
	"github.com/stretchr/testify/assert"
)

func TestGetWeek(t *testing.T) {
	expected := []time.Time{
		time.Date(2023, 12, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 4, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 6, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 7, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 8, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 9, 0, 0, 0, 0, time.UTC),
	}

	day := time.Date(2023, 12, 5, 0, 0, 0, 0, time.UTC)

	actual := timeutil.GetWeek(day)

	assert.Equal(t, expected, actual)
}
