package daily

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/apognu/gocal"
	"github.com/pkg/errors"
	"github.com/spa-stc/newsletter/server/profile"
	"github.com/spa-stc/newsletter/store"
	"github.com/spa-stc/newsletter/timeutil"
)

func Run(p *profile.Profile, s *store.Store) {
	slog.Info("beginning daily cron job")

	week := timeutil.GetWeek(time.Now())
	days := make([]store.Day, 7)
	for _, date := range week {
		day := store.Day{
			Date: date.Format(store.DayFormat),
		}

		days = append(days, day)
	}

	cal, err := GetGocal("")
	if err != nil {
		slog.Error("failed to parse ical", "error", err)
	}

	for _, day := range days {
		GetLunch(cal, day)
	}

	slog.Info("daily cron job completed")
}

func GetLunch(parser *gocal.Gocal, day store.Day) {
	for _, event := range parser.Events {
		if event.Start.UTC().Format(store.DayFormat) == day.Date {
			day.Lunch = event.Description
		}
	}

	day.Lunch = "Not Available"
}

func GetGocal(url string) (*gocal.Gocal, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query sheet url")
	}
	defer res.Body.Close()

	cal := gocal.NewParser(res.Body)
	cal.SkipBounds = true

	err = cal.Parse()
	if err != nil {
		return nil, err
	}

	return cal, nil
}
