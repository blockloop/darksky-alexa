package speakers

import (
	"fmt"
	"time"

	timeclock "github.com/benbjohnson/clock"
	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/darksky-alexa/pollen"
	"github.com/blockloop/darksky-alexa/tz"
	nowutil "github.com/jinzhu/now"
)

var clock = timeclock.New()

// NoData means there was no forecast data for the specified time period
var NoData = "I don't have any forecast information for that time"

type Speaker interface {
	CanSpeak(*alexa.WeatherRequest) bool
	Speak(*tz.Location, *darksky.Forecast, *pollen.Forecast, *alexa.WeatherRequest) string
	Name() string
}

// All is the ordered list of speakers
var All = []Speaker{
	Pollen{},
	Precipitation{},
	LowHigh{},
	Alerts{},
	Forecast{},
}

// Weather conditions that have been configured in the Alexa skill
// these show up under alexa.WeatherRequest.Condition
const (
	dayFormat = "2006-01-02"

	condPollen      = "pollen"
	condLow         = "low"
	condHigh        = "high"
	condSnow        = "snow"
	condRain        = "rain"
	condAlerts      = "alerts"
	condHumidity    = "humidity"
	condForecast    = "forecast"
	condExtForecast = "extended forecast"
	condTemperature = "temperature"
)

func sameWeek(a, b time.Time) bool {
	ayear, aweek := a.ISOWeek()
	byear, bweek := b.ISOWeek()
	return ayear == byear && aweek == bweek
}

func sameDay(a, b time.Time, times ...time.Time) bool {
	bod := nowutil.New(a.UTC()).BeginningOfDay()
	if bod != nowutil.New(b.UTC()).BeginningOfDay() {
		return false
	}

	for _, t := range times {
		if bod != nowutil.New(t.UTC()).BeginningOfDay() {
			return false
		}
	}
	return true
}

func today(b time.Time) bool {
	return sameDay(clock.Now(), b)
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

func humanPollen(index float64) string {
	if index < 2.5 {
		return "low"
	}
	if index < 4.9 {
		return "low-medium"
	}
	if index < 7.2 {
		return "medium"
	}
	if index < 9.7 {
		return "medium-high"
	}
	return "high"
}
