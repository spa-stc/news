package daysfetch

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/apognu/gocal"
	"github.com/gocarina/gocsv"
)

type DataGetter interface {
	GetLunch(ctx context.Context) (map[string]string, error)
	GetInfo(ctx context.Context) (map[string]CsvData, error)
}

type Getter struct {
	sheetID  string
	sheetGID string
	lunchURL string
}

func NewGetter(sheetID string, sheetGID string, lunchURL string) *Getter {
	return &Getter{
		sheetID:  sheetID,
		sheetGID: sheetGID,
		lunchURL: lunchURL,
	}
}

func (i *Getter) GetLunch(ctx context.Context) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, i.lunchURL, nil)
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

const CSVDateFormat = "1/2/2006"

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

func (i *Getter) GetInfo(ctx context.Context) (map[string]CsvData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	queryurl := fmt.Sprintf(
		"https://docs.google.com/spreadsheets/d/1EH8eAXgtaCzxBQWmk1yjBIuKD_CtZlsUDimHUFyoZs8/export?format=csv&id=%s&gid=%s",
		i.sheetID, i.sheetGID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryurl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}
	defer res.Body.Close()

	var data []CsvData
	if err := gocsv.Unmarshal(res.Body, &data); err != nil {
		return nil, err
	}

	results := make(map[string]CsvData)
	for _, v := range data {
		results[v.Date] = v
	}

	return results, nil
}
