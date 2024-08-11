package api

import (
	"net/http"
	"time"

	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/resource"
	"stpaulacademy.tech/newsletter/util/service"
	"stpaulacademy.tech/newsletter/util/timeutil"
	"stpaulacademy.tech/newsletter/web"
)

type WeekResponse struct {
	Days []resource.Day `json:"days"`
}

func handleGetWeek(db db.Executor, timegetter service.TimeGenerator) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		weekdates := timeutil.GetWeek(time.Sunday, timegetter.NowUTC())

		days, err := resource.GetManyDays(r.Context(), db, weekdates)
		if err != nil {
			return err
		}

		return web.Encode(w, http.StatusOK, &WeekResponse{
			Days: days,
		})
	}
}
