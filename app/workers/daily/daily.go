package daily

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/apognu/gocal"
	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
	"github.com/spa-stc/newsletter/server/profile"
	"github.com/spa-stc/newsletter/store"
	"github.com/spa-stc/newsletter/timeutil"
)

func Run(p *profile.Profile, s *store.Store) {
	ctx := context.Background()
	slog.Info("beginning daily cron job")

	week := timeutil.GetWeek(time.Now())
	var days []store.Day
	for _, date := range week {
		day := store.Day{
			Date: date.Format(store.DayFormat),
		}

		days = append(days, day)
	}

	cal, err := GetGocal(p.ICALURL)
	if err != nil {
		slog.Error("failed to parse ical", "error", err)
		return
	}

	csvdata, err := GetCsvData(p.SheedID, p.SheetName)
	if err != nil {
		slog.Error("error fetching csv data", "error", err)
		return
	}

	for _, day := range days {
		day := GetLunch(cal, day)

		day, err := GetDay(csvdata, day)
		if err != nil {
			slog.Error("error getting day from csv data", "error", err)
			return
		}

		_, err = s.UpsertDay(ctx, day)
		if err != nil {
			slog.Error("error inserting day into database", "error", err)
			return
		}
	}

	slog.Info("daily cron job completed")
}

func GetLunch(parser *gocal.Gocal, day store.Day) store.Day {
	for _, event := range parser.Events {
		if event.Start.UTC().Format(store.DayFormat) == day.Date {
			day.Lunch = event.Description
		}
	}

	day.Lunch = "Not Available"

	return day
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

type CsvData struct {
	Date     string `csv:"DATE"`
	Rday     string `csv:"R. DAY"`
	Location string `csv:"LOCATION"`
	Event    string `csv:"EVENT"`
	Grade9   string `csv:"9th GRADE"`
	Grade10  string `csv:"10th GRADE"`
	Grade11  string `csv:"11th GRADE"`
	Grade12  string `csv:"12th GRADE"`
	ApInfo   string `csv:"AP EXAMS"`
	CcInfo   string `csv:"CC TOPICS"`
}

func GetCsvData(id, sheet string) ([]CsvData, error) {
	queryurl := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv&sheet=%s", id, sheet)

	req, err := http.Get(queryurl)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching csv")
	}
	defer req.Body.Close()

	var data []CsvData
	if err := gocsv.Unmarshal(req.Body, &data); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling csv")
	}

	return data, nil
}

func GetDay(csvdata []CsvData, day store.Day) (store.Day, error) {
	time, err := time.Parse(store.DayFormat, day.Date)
	if err != nil {
		return store.Day{}, errors.Wrap(err, "error parsing time")
	}

	templ := time.Format("1/2/2006")
	for _, v := range csvdata {
		if v.Date == templ {
			day.RotationDay = v.Rday
			day.Location = v.Location
			day.XPeriod = v.Event
			day.Grade9 = v.Grade9
			day.Grade10 = v.Grade10
			day.Grade11 = v.Grade11
			day.Grade12 = v.Grade12
			day.ApInfo = v.ApInfo
			day.CCInfo = v.CcInfo
		}
	}

	return day, nil
}
