package api

import (
	"net/http"
	"time"

	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/resource"
	"stpaulacademy.tech/newsletter/util/service"
	"stpaulacademy.tech/newsletter/util/sliceutil"
	"stpaulacademy.tech/newsletter/util/timeutil"
	"stpaulacademy.tech/newsletter/web"
)

type Day struct {
	Date        time.Time `json:"date"`
	Lunch       string    `json:"lunch"`
	XPeriod     string    `json:"x_period"`
	RotationDay string    `json:"r_day"`
	Location    string    `json:"location"`
	Notes       string    `json:"notes"`
	ApInfo      string    `json:"ap_info"`
	CcInfo      string    `json:"cc_info"`
	Grade9      string    `json:"grade_9"`
	Grade10     string    `json:"grade_10"`
	Grade11     string    `json:"grade_11"`
	Grade12     string    `json:"grade_12"`
	CreatedTS   time.Time `json:"created_ts"`
	UpdatedTS   time.Time `json:"updated_ts"`
}

type WeekResponse struct {
	Days []Day `json:"days"`
}

func handleGetWeek(db db.Executor, timegetter service.TimeGenerator) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		weekdates := timeutil.GetWeek(time.Sunday, timegetter.NowUTC())

		days, err := resource.GetManyDays(r.Context(), db, weekdates)
		if err != nil {
			return err
		}

		d := sliceutil.Map(days, func(d resource.Day) Day {
			return Day{
				Date:        d.Date,
				Lunch:       d.Lunch,
				XPeriod:     d.XPeriod,
				RotationDay: d.RotationDay,
				Location:    d.Location,
				Notes:       d.Notes,
				ApInfo:      d.ApInfo,
				CcInfo:      d.CcInfo,
				Grade9:      d.Grade9,
				Grade10:     d.Grade10,
				Grade11:     d.Grade11,
				Grade12:     d.Grade12,
				CreatedTS:   d.CreatedTS,
				UpdatedTS:   d.UpdatedTS,
			}
		})

		return web.Encode(w, http.StatusOK, &WeekResponse{
			Days: d,
		})
	}
}
