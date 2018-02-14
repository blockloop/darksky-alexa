package speech

import (
	"fmt"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
)

// Weather conditions that have been configured in the Alexa skill
// these show up under alexa.WeatherRequest.Condition
const (
	Low         = "low"
	High        = "high"
	Snow        = "snow"
	Rain        = "rain"
	Humidity    = "humidity"
	Forecast    = "forecast"
	ExtForecast = "extended forecast"
	Temperature = "temperature"
)

// Speak formats weather that makes sense to the query
//
// if the question is "will it rain tomorrow?" then the response should
// respond with yes or no and the
func Speak(forecast darksky.Forecast, q alexa.WeatherRequest) string {
	ll := log.WithFields(log.Fields{
		"query.condition": q.Condition,
		"query.span":      q.Span,
		"query.start":     q.Start,
	})

	for _, speaker := range speakers {
		if msg, ok := speaker(&forecast, &q); ok {
			return msg
		}
	}

	ll.Info("using simple")
	msg := fmt.Sprintf("currently %.0f° and %s with a high of %.0f and a low of %.0f",
		forecast.Hourly.Data[0].Temperature,
		forecast.Currently.Summary,
		forecast.Daily.Data[0].TemperatureHigh,
		forecast.Daily.Data[0].TemperatureLow,
	)
	return msg
}
