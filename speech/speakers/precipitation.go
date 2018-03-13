package speakers

import (
	"fmt"
	"time"

	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/darksky-alexa/pollen"
)

// Precipitation responds to the following queries
// will it rain|snow [day] [time]
type Precipitation struct{}

func (Precipitation) Name() string {
	return "Precipitation"
}

func (Precipitation) CanSpeak(q *alexa.WeatherRequest) bool {
	return q.Condition == condRain || q.Condition == condSnow
}

func (p Precipitation) Speak(f *darksky.Forecast, _ *pollen.Forecast, q *alexa.WeatherRequest) string {
	dps := p.dataPoints(&q.TimeSpan, f)
	if len(dps) == 0 {
		return NoData
	}

	var found *darksky.DataPoint
	for _, data := range dps {
		if data.PrecipType != q.Condition || data.PrecipProbability < 0.10 {
			continue
		}
		if found == nil || data.PrecipProbability > found.PrecipProbability {
			found = &data
		}
	}

	if found == nil {
		return fmt.Sprintf("There is no chance of %s", q.Condition)
	}

	return fmt.Sprintf("There is %s of %s",
		humanProbability(found.PrecipProbability),
		q.Condition)
}

func (Precipitation) dataPoints(ts *alexa.TimeSpan, f *darksky.Forecast) []darksky.DataPoint {
	var prospects *darksky.DataBlock

	start := ts.Start.Add(-time.Nanosecond)
	end := ts.End.Add(time.Nanosecond)
	if today(start) && today(end) {
		prospects = &f.Hourly
	} else {
		prospects = &f.Daily
	}

	items := darksky.Where(prospects.Data, func(dp darksky.DataPoint) bool {
		itemTime := dp.Time.Time()
		return itemTime.After(start) && itemTime.Before(end)
	})
	return items
}
