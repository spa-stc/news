package daysfetch

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/apognu/gocal"
)

// Get lunch data from any source (the web or testing).
type LunchGetter interface {
	Get(ctx context.Context) (map[string]string, error)
}

type IcalLunchGetter struct {
	url string
}

func NewIcalLunchGetter(url string) *IcalLunchGetter {
	return &IcalLunchGetter{
		url: url,
	}
}

func (i *IcalLunchGetter) Get(ctx context.Context) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, i.url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
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
