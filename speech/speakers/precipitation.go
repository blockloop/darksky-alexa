package speakers

import (
	"fmt"

	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
)

// Precipitation responds to the following queries
// will it rain|snow [day] [time]
type Precipitation struct{}

func (Precipitation) Name() string {
	return "Precipitation"
}

func (Precipitation) CanSpeak(q *alexa.WeatherRequest) bool {
	return q.Condition == rain || q.Condition == snow
}

func (Precipitation) Speak(f *darksky.Forecast, q *alexa.WeatherRequest) string {
	dps := dataPoints(&q.TimeSpan, f)
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
