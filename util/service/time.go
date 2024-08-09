package service

import "time"

type TimeGenerator interface {
	// NowUTC returns the current time. This may be a stubbed time if the time
	// has been actively stubbed in a test.
	NowUTC() time.Time
}

type TimeGen struct{}

func NewTimeGen() *TimeGen {
	return &TimeGen{}
}

func (t *TimeGen) NowUTC() time.Time {
	return time.Now().UTC()
}
