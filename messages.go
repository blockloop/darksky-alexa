package main

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/go-chi/render"
)

func ResponseText(msg string) render.M {
	return render.M{
		"response": render.M{
			"outputSpeech": render.M{
				"type": "PlainText",
				"text": msg,
			},
			"shouldEndSession": true,
		},
		"sessionAttributes": render.M{},
	}
}

type Request struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	Intent    struct {
		Name  string `json:"name"`
		Slots struct {
			Condition struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"condition"`
			City struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"city"`
			Time struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"time"`
			Day struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"day"`
		} `json:"slots"`
	} `json:"intent"`
	Locale    string    `json:"locale"`
	Timestamp time.Time `json:"timestamp"`
}

type RequestValues struct {
	Condition string
	City      string
	Time      time.Time
	Duration  time.Duration
}

func (r Request) Values() RequestValues {
	slots := r.Intent.Slots
	res := RequestValues{
		Condition: slots.Condition.Value,
		City:      slots.City.Value,
	}

	timeofday := parseTime(slots.Time.Value)
	day, dur, err := parseDay(slots.Day.Value)
	if err != nil {
		log.WithError(err).Info("failed to parse day utterance")
	}
	res.Duration = dur

	timestamp := fmt.Sprintf("%s %s", day, timeofday)
	res.Time, err = time.Parse("2006-01-02 03:04", timestamp)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"timestamp": timestamp,
			"time":      timeofday,
			"day":       day,
		}).Error("failed to parse timestamp")
		res.Time = time.Now()
	}
	return res
}

func parseTime(s string) string {
	switch s {
	case "MO":
		return "08:00"
	default:
		return s
	}
}

func parseDay(s string) (day string, dur time.Duration, err error) {
	now := time.Now().Format("2006-01-02")
	// if !strings.Contains("-", s) {
	// 	return now, 1, fmt.Errorf("ignoring date utterance %q", s)
	// }
	// if !strings.Contains("W", s) {
	// 	return now, 1, nil
	// }

	// splits := strings.Split(s, "-")
	// // year := splits[0]
	// // date := time.Parse("2006", year)
	// if len(splits) == 3 {
	// 	// weekend
	// 	span = 2
	// } else {
	// 	span = 7
	// }

	// week, err := strconv.Atoi(strings.Replace(splits[1], "W", "", 1))
	// if err != nil {
	// 	return now, 1, fmt.Errorf("could not parse week number: %q", splits[1])
	// }

	return now, time.Hour * 24, nil
}
