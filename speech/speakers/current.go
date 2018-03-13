package speakers

import (
	"fmt"

	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/darksky-alexa/pollen"
)

// Current handles any request and responds with the current weather
// information. It is responsible for mainly the following requests
// but handles any request it is given
//
// the weather
// the current weather
type Current struct{}

func (Current) Name() string {
	return "Current"
}

func (Current) CanSpeak(q *alexa.WeatherRequest) bool {
	return true
}

func (Current) Speak(f *darksky.Forecast, pol *pollen.Forecast, _ *alexa.WeatherRequest) string {
	var pc pollen.DataPoint
	for _, point := range pol.DataPoints {
		if today(point.Day) {
			pc = point
			break
		}
	}

	return fmt.Sprintf("It's currently %.0f° and %s with a high of %.0f and a low of %.0f. Pollen is %s",
		f.Hourly.Data[0].Temperature,
		f.Currently.Summary,
		f.Daily.Data[0].TemperatureHigh,
		f.Daily.Data[0].TemperatureLow,
		humanPollen(pc.Index),
	)
}
