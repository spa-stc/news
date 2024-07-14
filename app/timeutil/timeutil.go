package timeutil

import "time"

// Returns the calendar week related with the inputted time.
func GetWeek(today time.Time) []time.Time {
	// Round to the nearest day.
	today = today.Round(86400000000000)

	weekday := int(today.Weekday())
	sunday := today.AddDate(0, 0, -weekday)

	var times []time.Time
	for v := 0; v < 7; v++ {
		times = append(times, sunday.AddDate(0, 0, v))
	}

	return times
}

func GetWeekday(today time.Time) string {
	return today.Format("Monday")
}
