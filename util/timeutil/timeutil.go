package timeutil

import "time"

const day = time.Duration(86400000000000) // day as a time.Duration.

// Get days of this week, with hour = 0, indexed from the given start weekday.
func GetWeek(start time.Weekday, today time.Time) []time.Time {
	today = today.Truncate(day)
	// Actual date coresponding to start.
	weekstart := today.AddDate(0, 0, -(int(today.Weekday() - start)))

	week := make([]time.Time, 7)
	for v := range 7 {
		week[v] = weekstart.AddDate(0, 0, v)
	}

	return week
}
