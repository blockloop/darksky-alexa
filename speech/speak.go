package speech

import (
	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/darksky-alexa/speech/speakers"
)

// Default speaks the default message
var Default = speakers.Current{}.Speak

// Speak formats weather that makes sense to the query. If no handler is
// found then the Default speaker is used.
//
// if the question is "will it rain tomorrow?" then the response should
// respond with yes or no and the time it is expected to rain
func Speak(forecast *darksky.Forecast, q alexa.WeatherRequest) string {
	ll := log.WithFields(log.Fields{
		"query.condition": q.Condition,
		"query.end":       q.End,
		"query.start":     q.Start,
	})

	for _, speaker := range speakers.All {
		if !speaker.CanSpeak(&q) {
			continue
		}
		ll.WithField("speaker", speaker.Name()).Info("speaking")
		return speaker.Speak(forecast, &q)
	}

	ll.Warn("No speaker was found. Falling back to default.")
	return Default(forecast, &q)
}
