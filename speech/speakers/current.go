package speakers

import (
	"fmt"

	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
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

func (Current) Speak(f *darksky.Forecast, _ *alexa.WeatherRequest) string {
	return fmt.Sprintf("It's currently %.0f° and %s with a high of %.0f and a low of %.0f",
		f.Hourly.Data[0].Temperature,
		f.Currently.Summary,
		f.Daily.Data[0].TemperatureHigh,
		f.Daily.Data[0].TemperatureLow,
	)
}
