package speakers

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
)

// LowHigh responds to the following queries
//
// the low|high [day] [time]
// the low|high
type LowHigh struct{}

func (LowHigh) Name() string {
	return "LowHigh"
}

func (lh LowHigh) CanSpeak(q *alexa.WeatherRequest) bool {
	return q.Condition == high || q.Condition == low
}

func (lh LowHigh) Speak(f *darksky.Forecast, q *alexa.WeatherRequest) string {
	if !lh.CanSpeak(q) {
		log.Error("tried to speak low/high without asking for low/high")
		return "a problem occurred"
	}

	ts := &q.TimeSpan
	dps := darksky.Where(f.Daily.Data, func(dp darksky.DataPoint) bool {
		itemTime := dp.Time.Time()
		return itemTime.Before(ts.End) && (itemTime.After(ts.Start) || sameDay(itemTime, ts.Start))
	})
	if len(dps) == 0 {
		return NoData
	}

	var temp float64
	var when time.Time
	switch q.Condition {
	case low:
		temp, when = findLow(dps)
		break
	case high:
		temp, when = findHigh(dps)
		break
	default:
		log.Error("tried to speak low/high without asking for low/high")
		return "a problem occurred"
	}

	return fmt.Sprintf("%.0f is the %s for %s", temp, q.Condition, humanDay(when))
}

func findLow(dps []darksky.DataPoint) (lowest float64, when time.Time) {
	lowest = dps[0].TemperatureLow
	when = dps[0].Time.Time()

	for _, dp := range dps {
		if dp.TemperatureLow < lowest {
			lowest = dp.TemperatureLow
			when = dp.Time.Time()
		}
	}
	return lowest, when
}

func findHigh(dps []darksky.DataPoint) (high float64, when time.Time) {
	high = dps[0].TemperatureHigh
	when = dps[0].Time.Time()

	for _, dp := range dps {
		if dp.TemperatureHigh > high {
			high = dp.TemperatureHigh
			when = dp.Time.Time()
		}
	}
	return high, when
}
