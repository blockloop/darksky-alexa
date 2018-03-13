package speakers

import (
	"fmt"

	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/darksky-alexa/pollen"
)

// Pollen responds to the following queries
// the pollen for [day]
type Pollen struct{}

func (Pollen) Name() string {
	return "Pollen"
}

func (Pollen) CanSpeak(q *alexa.WeatherRequest) bool {
	return q.Condition == condPollen
}

func (p Pollen) Speak(_ *darksky.Forecast, pf *pollen.Forecast, q *alexa.WeatherRequest) string {
	for _, dp := range pf.DataPoints {
		if sameDay(dp.Day, q.Start) {
			return fmt.Sprintf("Pollen is %s %s",
				humanPollen(dp.Index),
				humanDay(dp.Day))
		}
	}
	return NoData
}
