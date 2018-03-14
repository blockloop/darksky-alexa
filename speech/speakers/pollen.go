package speakers

import (
	"fmt"

	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/darksky-alexa/pollen"
	"github.com/blockloop/darksky-alexa/tz"
)

// Pollen responds to the following queries
// the pollen for [day]
type Pollen struct{}

// Name is the name of this speaker
func (Pollen) Name() string {
	return "Pollen"
}

// CanSpeak returns true if this speaker can speak for this query
func (Pollen) CanSpeak(q *alexa.WeatherRequest) bool {
	return q.Condition == condPollen
}

// Speak speaks for this query
func (p Pollen) Speak(loc *tz.Location, _ *darksky.Forecast, pf *pollen.Forecast, q *alexa.WeatherRequest) string {
	for _, dp := range pf.DataPoints {
		if sameDay(dp.Day, q.Start) {
			return fmt.Sprintf("Pollen is %s %s",
				humanPollen(dp.Index),
				humanDay(dp.Day))
		}
	}
	return NoData
}
