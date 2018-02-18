package speakers

import (
	"fmt"
	"sort"

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

	switch q.Condition {
	case low:
		dp := findLow(dps)
		return fmt.Sprintf("%.0f is the %s %s", dp.TemperatureLow, q.Condition, humanDay(dp.Time.Time()))
	case high:
		dp := findHigh(dps)
		return fmt.Sprintf("%.0f is the %s %s", dp.TemperatureHigh, q.Condition, humanDay(dp.Time.Time()))
	default:
		log.Error("tried to speak low/high without asking for low/high")
		return "a problem occurred"
	}
}

func findLow(dps []darksky.DataPoint) *darksky.DataPoint {
	sort.SliceStable(dps, func(i, j int) bool {
		return dps[i].TemperatureLow < dps[j].TemperatureLow
	})
	return &dps[0]
}

func findHigh(dps []darksky.DataPoint) *darksky.DataPoint {
	sort.SliceStable(dps, func(i, j int) bool {
		return dps[i].TemperatureHigh > dps[j].TemperatureHigh
	})
	return &dps[0]
}
