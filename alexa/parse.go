package alexa

import (
	"regexp"
	"strconv"
	"time"

	"github.com/apex/log"
	nowutil "github.com/jinzhu/now"
)

// WeatherRequest is a formatted RequestBody
type WeatherRequest struct {
	Condition string
	TimeSpan
}

var (
	weekRegex = regexp.MustCompile(`(\d{4})-W(\d{1,2})(-WE)?`)
)

const (
	// Day is a time.Duration of 24 hours
	Day time.Duration = time.Hour * 24
	// Week is a time.Duration of 7 days
	Week time.Duration = Day * 7
)

const (
	dayFmt   = "2006-01-02"
	monthFmt = "2006-01"

	morning   = "MO"
	afternoon = "AF"
	evening   = "EV"
	night     = "NI"
	weekend   = "WE"
)

// ParseWeatherRequest parses the RequestBody as a WeatherRequest
func ParseWeatherRequest(r RequestBody) WeatherRequest {
	slots := r.Intent.Slots
	res := WeatherRequest{
		Condition: slots.Condition.Value,
	}

	res.TimeSpan = parseDuration(slots.Time.Value, slots.Day.Value)
	return res
}

func parseDuration(timeofday, day string) TimeSpan {
	ts := parseDay(day)
	if timeofday != "" {
		hour, span := parseTime(timeofday)
		ts.Start = time.Date(ts.Start.Year(), ts.Start.Month(), ts.Start.Day(), hour, 0, 0, 0, time.Local)
		ts.End = ts.Start.Add(span)
	}
	return ts
}

// parseTime parses an AMAZON.TIME
//
// night: NI, morning: MO, afternoon: AF, evening: EV. or HH:MM
func parseTime(s string) (hour int, dur time.Duration) {
	switch s {
	case morning:
		return 6, time.Hour * 6
	case afternoon:
		return 12, time.Hour * 5
	case evening:
		return 17, time.Hour * 4
	case night:
		return 20, time.Hour * 4
	}

	parsed, err := time.Parse("15:04", s)
	if err != nil {
		log.WithField("time", s).Warn("invalid time")
		return 0, Day
	}
	return parsed.Hour(), time.Hour
}

// parseDay parses an AMAZON.DATE
//
// Utterances that map to a specific date (such as “today”, “now”, or “november twenty-fifth”)
// convert to a complete date: 2015-11-25. Note that this defaults to dates on or after
// the current date
//
// Utterances that map to just a specific week (such as “this week” or “next week”),
// convert a date indicating the week number: 2015-W49.
//
// Utterances that map to the weekend for a specific week (such as “this weekend”)
// convert to a date indicating the week number and weekend: 2015-W49-WE.
//
// Utterances that map to a month, but not a specific day (such as “next month”, or “december”)
// convert to a date with just the year and month: 2015-12.
//
// Utterances that map to a year (such as “next year”) convert to a date containing just the year: 2016.
//
// Utterances that map to a decade convert to a date indicating the decade: 201X.
// Utterances that map to a season (such as “next winter”) convert to a date with the year
// and a season indicator: winter: WI, spring: SP, summer: SU, fall: FA)
//
// https://developer.amazon.com/docs/custom-skills/slot-type-reference.html#date
func parseDay(day string) (ts TimeSpan) {
	tnow := time.Now()
	now := TimeSpan{
		Start: tnow,
		End:   tnow,
	}

	if day == "" {
		return now
	}

	// try to parse simple date first
	if simpleDay, err := time.Parse(dayFmt, day); err == nil {
		return TimeSpan{
			Start: simpleDay,
			End:   nowutil.New(simpleDay).EndOfDay(),
		}
	}
	if simpleMonth, err := time.Parse(monthFmt, day); err == nil {
		return TimeSpan{
			Start: simpleMonth,
			End:   nowutil.New(simpleMonth).EndOfWeek(),
		}
	}

	if week := parseWeek(day); !week.IsZero() {
		return week
	}

	log.WithField("day", day).Info("unknown day. Using NOW")
	return now
}

func parseWeek(s string) TimeSpan {
	items := weekRegex.FindStringSubmatch(s)
	// no match
	if len(items) == 0 || items[0] == "" {
		return TimeSpan{}
	}

	year, _ := strconv.Atoi(items[1])
	zweek, _ := strconv.Atoi(items[2])
	zweek = zweek - 1
	weekend := items[3] != ""

	week := nowutil.New(time.Date(year, time.January, 1, 0, 0, 0, 1, time.Local))
	if zweek > 0 {
		week = nowutil.New(week.AddDate(0, 0, zweek*7))
	}
	bow := week.BeginningOfWeek()
	eow := week.EndOfWeek()

	log.WithFields(log.Fields{
		"orig":   s,
		"parsed": week.String(),
	}).Info("parsed date")

	if !weekend {
		return TimeSpan{
			Start: bow,
			End:   eow,
		}
	}
	// weekend means Saturday and Sunday
	eow = eow.AddDate(0, 0, 1)
	return TimeSpan{
		Start: nowutil.New(eow.AddDate(0, 0, -1)).BeginningOfDay(),
		End:   eow,
	}
}
