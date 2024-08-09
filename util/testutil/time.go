package testutil

import "time"

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
