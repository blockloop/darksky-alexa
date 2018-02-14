package speech

import (
	"fmt"
	"strings"
	"time"

	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
)

type speaker func(*darksky.Forecast, *alexa.WeatherRequest) (string, bool)

// NoData means there was no forecast data for the specified time period
var NoData = "I don't have any forecast information for that time"

var speakers = []speaker{
	speakYesNo,
	speakLowHigh,
	speakForecast,
	speakNow,
}

// the weather
// the current weather
func speakNow(f *darksky.Forecast, _ *alexa.WeatherRequest) (string, bool) {
	msg := fmt.Sprintf("It's currently %.0f° and %s with a high of %.0f and a low of %.0f",
		f.Hourly.Data[0].Temperature,
		f.Currently.Summary,
		f.Daily.Data[0].TemperatureHigh,
		f.Daily.Data[0].TemperatureLow,
	)
	return msg, true
}

// will it rain|snow [day] [time]
func speakYesNo(f *darksky.Forecast, q *alexa.WeatherRequest) (string, bool) {
	if q.Condition != Rain && q.Condition != Snow {
		return "", false
	}
	dps := dataPoints(&q.TimeSpan, f)
	if len(dps) == 0 {
		return NoData, true
	}

	var found *darksky.DataPoint
	for _, data := range dps {
		if data.PrecipType != q.Condition || data.PrecipProbability < 10 {
			continue
		}
		if found == nil || data.PrecipProbability > found.PrecipProbability {
			found = data
		}
	}

	if found == nil {
		msg := fmt.Sprintf("I see no of %s in the forecast for that time period", q.Condition)
		return msg, true
	}

	// try to take the easy way out
	if strings.Contains(found.Summary, q.Condition) {
		return found.Summary, true
	}
	when := time.Unix(found.UnixTime, 0)
	msg := fmt.Sprintf("There is %s of %s on %s the %s at %s",
		humanProbability(found.PrecipProbability),
		found.PrecipType,
		when.Weekday(),
		when.Day(),
		when.Format(time.Kitchen))
	return msg, true
}

// the low|high [day] [time]
// the low|high
func speakLowHigh(f *darksky.Forecast, q *alexa.WeatherRequest) (string, bool) {
	switch q.Condition {
	case Low, High:
		break
	default:
		return "", false
	}

	start := q.Start.UTC()
	end := start.Add(q.Span)
	data := make([]*darksky.DataPoint, 0, len(f.Daily.Data))
	for _, dp := range f.Daily.Data {
		itemTime := time.Unix(dp.UnixTime, 0).UTC()
		if itemTime.Before(end) && itemTime.After(start) {
			data = append(data, &dp)
		}
	}
	if len(data) == 0 {
		return NoData, true
	}

	var value *float64
	var day time.Time
	for _, dp := range data {
		switch q.Condition {
		case Low:
			if value == nil || dp.TemperatureLow < *value {
				value = &dp.TemperatureLow
				day = time.Unix(dp.UnixTime, 0)
			}
		case High:
			if value == nil || dp.TemperatureHigh > *value {
				value = &dp.TemperatureHigh
				day = time.Unix(dp.UnixTime, 0)
			}
		}
	}

	today := time.Now().UTC()
	verb := "is"
	if today.After(day.UTC()) {
		verb = "was"
	}

	msg := fmt.Sprintf("The %s %s %.0f", q.Condition, verb, *value)
	if !sameDay(today, day.UTC()) {
		msg = fmt.Sprintf("%s on %s", msg, day.Weekday())
	}
	return msg, true
}

// the weather|forecast [day] [time]
func speakForecast(*darksky.Forecast, *alexa.WeatherRequest) (string, bool) {
	return "", false
}

// dataPoints returns the data points that pertain to the requested time
func dataPoints(ts *alexa.TimeSpan, f *darksky.Forecast) []*darksky.DataPoint {
	var prospects darksky.DataPoints

	if ts.Span.Hours() > 24 {
		prospects = f.Daily
	} else if ts.Span.Hours() >= 1 {
		prospects = f.Hourly
	} else {
		prospects = f.Minutely
	}

	start := ts.Start.Add(-time.Nanosecond)
	end := ts.Start.Add(ts.Span)

	data := make([]*darksky.DataPoint, 0, len(prospects.Data))
	for _, dp := range prospects.Data {
		eventTime := time.Unix(dp.UnixTime, 0).UTC()
		if eventTime.After(start) && eventTime.Before(end) {
			data = append(data, &dp)
		}
	}

	return data
}

func humanProbability(percent float64) string {
	if percent > 60 {
		return "a great chance"
	}
	if percent > 25 {
		return "a good chance"
	}
	if percent > 9 {
		return "a slight chance"
	}
	return "no chance"
}

func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

func sameHour(a, b time.Time) bool {
	return sameDay(a, b) && a.Hour() == b.Hour()
}

func sameMinute(a, b time.Time) bool {
	return sameHour(a, b) && a.Minute() == b.Minute()
}

func maxInt(a, b int) int {
	if a == b || a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a == b || a < b {
		return a
	}
	return b
}
