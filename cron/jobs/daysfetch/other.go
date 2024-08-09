package daysfetch

import (
	"fmt"
	"net/http"

	"github.com/gocarina/gocsv"
)

const CSVDateFormat = "1/2/2006"

type OtherInfoGetter interface {
	Get() (map[string]CsvData, error)
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

type OtherInfoCSVGetter struct {
	sheetID   string
	sheetName string
}

func NewCsvGetter(sheetID string, sheetName string) *OtherInfoCSVGetter {
	return &OtherInfoCSVGetter{
		sheetID:   sheetID,
		sheetName: sheetName,
	}
}

func (s *OtherInfoCSVGetter) Get() (map[string]CsvData, error) {
	queryurl := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv&sheet=%s", s.sheetID, s.sheetName)

	req, err := http.Get(queryurl)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	var data []CsvData
	if err := gocsv.Unmarshal(req.Body, &data); err != nil {
		return nil, err
	}

	results := make(map[string]CsvData)
	for _, v := range data {
		results[v.Date] = v
	}

	return results, nil
}
