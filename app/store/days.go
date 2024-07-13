package store

const DayFormat = "2006-01-02"

// Main Day Model, Outputted From Cron Jobs.
type Day struct {
	// Day format, 2006-01-02.
	Date,

	Lunch,
	XPeriod,
	RotationDay,
	Location,
	Notes,
	ApInfo,
	CCInfo,
	Grade9,
	Grade10,
	Grade11,
	Grade12 string

	CreatedTs int64
	UpdatedTs int64
}

type FindDay struct {
	Date string
}

type FindDays struct {
	Dates []string
}
