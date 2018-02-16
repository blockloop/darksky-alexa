package speakers

import (
	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
)

type Forecast struct{}

func (Forecast) Name() string {
	return "Forecast"
}

func (Forecast) CanSpeak(*alexa.WeatherRequest) bool {
	return false
}

// the weather|forecast [day] [time]
func (Forecast) Speak(*darksky.Forecast, *alexa.WeatherRequest) string {
	return ""
}
