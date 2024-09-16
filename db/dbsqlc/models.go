// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package dbsqlc

import (
	"time"
)

type Announcement struct {
	ID           int64
	Title        string
	Author       string
	Content      string
	DisplayStart time.Time
	DisplayEnd   time.Time
	CreatedTs    time.Time
	UpdatedTs    time.Time
}

type Day struct {
	Date        time.Time
	Lunch       string
	XPeriod     string
	RotationDay string
	Location    string
	Notes       string
	CcInfo      string
	Grade9      string
	Grade10     string
	Grade11     string
	Grade12     string
	CreatedTs   time.Time
	UpdatedTs   time.Time
}
