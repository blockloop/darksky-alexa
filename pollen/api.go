package pollen

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

const apiurlFmt = `https://www.pollen.com/api/forecast/extended/pollen/%s`

// API interacts with the pollen count API
type API struct {
	client *http.Client
}

func NewAPI() *API {
	return &API{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (api *API) GetPollen(ctx context.Context, zipcode string) (*Forecast, error) {
	url := fmt.Sprintf(apiurlFmt, zipcode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate request")
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.146 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", url)

	resp, err := api.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "failed http request")
	}

	var res httpResp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	f := &Forecast{
		Location: Location{
			Zip:   zipcode,
			City:  res.Location.City,
			State: res.Location.State,
		},
		DataPoints: make([]DataPoint, len(res.Location.Periods)),
	}

	for i, p := range res.Location.Periods {
		period := p.Period
		date, err := time.Parse("2006-01-02", strings.Split(period, "T")[0])
		if err != nil {
			log.WithField("raw", period).
				WithError(err).
				Error("failed to parse date")
			continue
		}

		f.DataPoints[i] = DataPoint{
			Day:   date,
			Index: p.Index,
		}
	}

	return f, nil
}

type httpResp struct {
	Type         string `json:"Type"`
	ForecastDate string `json:"ForecastDate"`
	Location     struct {
		ZIP     string `json:"ZIP"`
		City    string `json:"City"`
		State   string `json:"State"`
		Periods []struct {
			Period string  `json:"Period"`
			Index  float64 `json:"Index"`
		} `json:"periods"`
		DisplayLocation string `json:"DisplayLocation"`
	} `json:"Location"`
}

type Forecast struct {
	Location   Location
	DataPoints []DataPoint
}

type Location struct {
	Zip   string
	City  string
	State string
}

type DataPoint struct {
	Day   time.Time
	Index float64
}
