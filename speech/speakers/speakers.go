package speakers

import (
	"fmt"
	"time"

	timeclock "github.com/benbjohnson/clock"
	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
)

var clock = timeclock.New()

// NoData means there was no forecast data for the specified time period
var NoData = "I don't have any forecast information for that time"

type Speaker interface {
	CanSpeak(*alexa.WeatherRequest) bool
	Speak(*darksky.Forecast, *alexa.WeatherRequest) string
	Name() string
}

// All is the ordered list of speakers
var All = []Speaker{Precipitation{}, LowHigh{}, Forecast{}, Current{}}

// Weather conditions that have been configured in the Alexa skill
// these show up under alexa.WeatherRequest.Condition
const (
	dayFormat = "2006-01-02"

	low         = "low"
	high        = "high"
	snow        = "snow"
	rain        = "rain"
	humidity    = "humidity"
	forecast    = "forecast"
	extForecast = "extended forecast"
	temperature = "temperature"
)

func sameWeek(a, b time.Time) bool {
	ayear, aweek := a.ISOWeek()
	byear, bweek := b.ISOWeek()
	return ayear == byear && aweek == bweek
}

func sameDay(a, b time.Time, times ...time.Time) bool {
	afmt, bfmt := a.Format(dayFormat), b.Format(dayFormat)

	if afmt != bfmt {
		return false
	}

	for _, t := range times {
		if t.Format(dayFormat) != afmt {
			return false
		}
	}
	return true
}

func today(b time.Time) bool {
	return sameDay(clock.Now(), b)
}

func sameHour(a, b time.Time) bool {
	return a.Format("2006-01-02 15") == b.Format("2006-01-02 15")
}

func humanDay(day time.Time) string {
	today := clock.Now()

	if sameDay(today, day) {
		return "today"
	}
	if sameDay(today.AddDate(0, 0, 1), day) {
		return fmt.Sprintf("tomorrow")
	}
	if sameWeek(today, day) {
		return day.Weekday().String()
	}
	return day.Format("Monday, Jan _2")

}

func timeOfDay(t time.Time) string {
	hour := t.Hour()
	if hour <= 11 {
		return "morning"
	}
	if hour <= 17 {
		return "afternoon"
	}
	if hour <= 20 {
		return "evening"
	}
	return "night"
}

func humanProbability(percent float64) string {
	if percent > 0.40 {
		return "a very good chance"
	}
	if percent > 0.20 {
		return "a good chance"
	}
	if percent > 0.09 {
		return "a slight chance"
	}
	return "no chance"
}
