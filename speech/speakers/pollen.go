package speakers

import (
	"fmt"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/darksky-alexa/pollen"
)

// Pollen responds to the following queries
// the pollen for [day] [time]
type Pollen struct{}

func (Pollen) Name() string {
	return "Pollen"
}

func (Pollen) CanSpeak(q *alexa.WeatherRequest) bool {
	return q.Condition == polln
}

func (p Pollen) Speak(_ *darksky.Forecast, pf *pollen.Forecast, q *alexa.WeatherRequest) string {
	for _, dp := range pf.DataPoints {
		log.WithFields(log.Fields{
			"dp.day":  dp.Day,
			"q.start": q.Start,
		}).Info("comparing")
		if sameDay(dp.Day, q.Start) {
			return fmt.Sprintf("Pollen is %s %s",
				humanPollen(dp.Index),
				humanDay(dp.Day))
		}
	}
	return NoData
}
