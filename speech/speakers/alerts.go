package speakers

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/darksky-alexa/pollen"
)

const noAlerts = "there are no active alerts"

type Alerts struct{}

func (Alerts) Name() string {
	return "Alerts"
}

func (Alerts) CanSpeak(q *alexa.WeatherRequest) bool {
	return q.Condition == condAlerts
}

// Speak speaks "the weather|forecast [day] [time]"
func (f Alerts) Speak(forecast *darksky.Forecast, _ *pollen.Forecast, q *alexa.WeatherRequest) string {
	if !f.CanSpeak(q) {
		log.Error("attempted to speak Alerts without checking CanSpeak first")
		return "a problem occured"
	}

	now := time.Now()
	if !sameDay(q.Start, now) {
		return NoData
	}

	if len(forecast.Alerts) == 0 {
		return noAlerts
	}

	active := make([]darksky.Alert, 0, len(forecast.Alerts))
	for _, alert := range forecast.Alerts {
		if !alert.Expires.Time().After(now) {
			continue
		}
		active = append(active, alert)
	}

	if len(active) == 0 {
		return noAlerts
	}

	sort.SliceStable(active, func(i, j int) bool {
		return active[i].Time.Time().After(active[j].Time.Time())
	})

	// using a hash because the same message occurs multiple times
	// with different timestamps
	hash := map[string]string{}
	for _, alert := range active {
		exp := alert.Expires.Time()
		msg := fmt.Sprintf("%s remains in effect until %s at %s.\n",
			alert.Title,
			humanDay(exp),
			exp.Format(time.Kitchen))

		if _, ok := hash[alert.Title]; !ok {
			hash[alert.Title] = msg
		}
	}

	sb := bytes.NewBuffer(nil)
	for _, msg := range hash {
		sb.WriteString(msg)
	}
	return sb.String()
}
