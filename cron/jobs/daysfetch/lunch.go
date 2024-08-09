package daysfetch

import (
	"net/http"
	"time"

	"github.com/apognu/gocal"
)

// Get lunch data from any source (the web or testing).
type LunchGetter interface {
	Get() (map[string]string, error)
}

type IcalLunchGetter struct {
	url string
}

func NewIcalLunchGetter(url string) *IcalLunchGetter {
	return &IcalLunchGetter{
		url: url,
	}
}

func (i *IcalLunchGetter) Get() (map[string]string, error) {
	res, err := http.Get(i.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	cal := gocal.NewParser(res.Body)

	cal.SkipBounds = true

	err = cal.Parse()
	if err != nil {
		return nil, err
	}

	lunchData := make(map[string]string)
	for _, v := range cal.Events {
		startString := v.Start.UTC().Format(time.DateOnly)

		lunchData[startString] = v.Description
	}

	return lunchData, nil
}
