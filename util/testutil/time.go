package testutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type TestTimeGen struct {
	now time.Time
}

func NewTestTimeGen(now time.Time) *TestTimeGen {
	return &TestTimeGen{
		now: now,
	}
}

func (t *TestTimeGen) NowUTC() time.Time {
	return t.now
}

func ParseDate(t *testing.T, s string) time.Time {
	date, err := time.Parse(time.DateOnly, s)
	require.NoError(t, err)
	return date
}
