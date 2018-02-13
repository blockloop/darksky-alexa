package alexa

import (
	"fmt"
	"time"

	"github.com/apex/log"
)

// JSON is json
type JSON map[string]interface{}

// ResponseText creates a response message to send to an alexa request
func ResponseText(msg string) interface{} {
	return JSON{
		"response": JSON{
			"outputSpeech": JSON{
				"type": "PlainText",
				"text": msg,
			},
			"shouldEndSession": true,
		},
		"sessionAttributes": JSON{},
	}
}

// Request is found here
// https://developer.amazon.com/docs/custom-skills/request-and-response-json-reference.html
type Request struct {
	Version string `json:"version"`
	Session struct {
		New         bool   `json:"new"`
		SessionID   string `json:"sessionId"`
		Application struct {
			ApplicationID string `json:"applicationId"`
		} `json:"application"`
		Attributes struct {
			Key string `json:"key"`
		} `json:"attributes"`
		User struct {
			UserID      string `json:"userId"`
			AccessToken string `json:"accessToken"`
			Permissions struct {
				ConsentToken string `json:"consentToken"`
			} `json:"permissions"`
		} `json:"user"`
	} `json:"session"`
	Context struct {
		System struct {
			Device struct {
				DeviceID            string `json:"deviceId"`
				SupportedInterfaces struct {
					AudioPlayer struct {
					} `json:"AudioPlayer"`
				} `json:"supportedInterfaces"`
			} `json:"device"`
			Application struct {
				ApplicationID string `json:"applicationId"`
			} `json:"application"`
			User struct {
				UserID      string `json:"userId"`
				AccessToken string `json:"accessToken"`
				Permissions struct {
					ConsentToken string `json:"consentToken"`
				} `json:"permissions"`
			} `json:"user"`
			APIEndpoint    string `json:"apiEndpoint"`
			APIAccessToken string `json:"apiAccessToken"`
		} `json:"System"`
		AudioPlayer struct {
			PlayerActivity       string `json:"playerActivity"`
			Token                string `json:"token"`
			OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
		} `json:"AudioPlayer"`
	} `json:"context"`
	Request RequestBody `json:"request"`
}

// RequestBody is the body of the specific request
type RequestBody struct {
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

// WeatherRequest is a formatted RequestBody
type WeatherRequest struct {
	Condition string
	City      string
	Time      time.Time
	Duration  time.Duration
}

// WeatherRequest parses the RequestBody as a WeatherRequest
func (r RequestBody) WeatherRequest() WeatherRequest {
	slots := r.Intent.Slots
	res := &WeatherRequest{
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
	return *res
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
