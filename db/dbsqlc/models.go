// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package dbsqlc

import (
	"time"
)

type Day struct {
	Date        time.Time
	Lunch       string
	XPeriod     string
	RotationDay string
	Location    string
	Notes       string
	ApInfo      string
	CcInfo      string
	Grade9      string
	Grade10     string
	Grade11     string
	Grade12     string
	CreatedTs   time.Time
	UpdatedTs   time.Time
}