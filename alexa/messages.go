package alexa

import (
	"time"
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
				Value string `json:"value"`
			} `json:"condition"`
			Time struct {
				Value string `json:"value"`
			} `json:"time"`
			Day struct {
				Value string `json:"value"`
			} `json:"day"`
		} `json:"slots"`
	} `json:"intent"`
	Locale    string    `json:"locale"`
	Timestamp time.Time `json:"timestamp"`
}
