package speakers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sort"
	"time"

	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"

	nowutil "github.com/jinzhu/now"
)

type Forecast struct{}

func (Forecast) Name() string {
	return "Forecast"
}

func (Forecast) CanSpeak(*alexa.WeatherRequest) bool {
	return true
}

// the weather|forecast [day] [time]
func (f Forecast) Speak(forecast *darksky.Forecast, q *alexa.WeatherRequest) string {
	now := time.Now()
	dur := q.End.Sub(q.Start)

	if dur.Hours() < 5 && q.End.Before(now.AddDate(0, 0, 2)) {
		return f.hourly(q.TimeSpan, forecast)
	}

	return f.daily(q.TimeSpan, forecast)

}

func (Forecast) daily(ts alexa.TimeSpan, forecast *darksky.Forecast) string {
	start := nowutil.New(ts.Start).BeginningOfDay().Add(-time.Nanosecond)
	end := nowutil.New(ts.End).EndOfDay()
	dps := darksky.Where(forecast.Daily.Data, func(dp darksky.DataPoint) bool {
		t := dp.Time.Time()
		return t.After(start) && t.Before(end)
	})

	switch len(dps) {
	case 0:
		return NoData
	case 1:
		dp := dps[0]
		return fmt.Sprintf("%s Temperatures will range between %.0f and %.0f",
			dp.Summary, dp.TemperatureLow, dp.TemperatureHigh)
	}

	buf := bytes.NewBufferString("")
	for i, dp := range dps {
		if i > 1 {
			break
		}
		fmt.Fprintf(buf, "On %s, %s The high will be %0.f and the low will be %0.f with %s of %s. ",
			humanDay(dp.Time.Time()),
			dp.Summary,
			dp.TemperatureHigh,
			dp.TemperatureLow,
			humanProbability(dp.PrecipProbability),
			dp.PrecipType)
	}
	b, _ := ioutil.ReadAll(buf)
	return string(b)
}

func (Forecast) hourly(ts alexa.TimeSpan, forecast *darksky.Forecast) string {
	start := nowutil.New(ts.Start).BeginningOfHour().Add(-time.Nanosecond)
	end := nowutil.New(ts.End).EndOfHour().Add(time.Nanosecond)
	dps := darksky.Where(forecast.Hourly.Data, func(dp darksky.DataPoint) bool {
		t := dp.Time.Time()
		return t.After(start) && t.Before(end)
	})

	if len(dps) == 0 {
		return NoData
	}

	avg, _ := avgTemp(dps)

	msg := fmt.Sprintf("About %0.f° and %s", avg, dps[0].Summary)

	if prec := highestChanceOfPrecipitation(dps); prec.PrecipProbability > 0.15 {
		msg = fmt.Sprintf(" with %s chance of %s",
			humanProbability(prec.PrecipProbability),
			prec.PrecipType)
	}
	return msg
}

func highestChanceOfPrecipitation(dps []darksky.DataPoint) *darksky.DataPoint {
	if len(dps) == 0 {
		return nil
	}
	if len(dps) == 1 {
		return &dps[0]
	}

	sort.Slice(dps, func(i, j int) bool {
		return dps[i].PrecipProbability < dps[j].PrecipProbability
	})

	return &dps[0]
}

func avgTemp(dps []darksky.DataPoint) (float64, bool) {
	if len(dps) == 0 {
		return 0, false
	}
	if len(dps) == 1 {
		return dps[0].Temperature, true
	}

	tot := 0.0
	for _, dp := range dps {
		tot += dp.Temperature
	}
	return tot / float64(len(dps)), true
}
