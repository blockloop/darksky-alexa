package handlers

import (
	"context"
	"net/http"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/darksky-alexa/geo"
	"github.com/blockloop/darksky-alexa/speech"
	"github.com/blockloop/tea"
	"github.com/pkg/errors"
)

// Ping is a simple handler that responds for healthchecks
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"pong": ""}`))
}

type darkskyAPI interface {
	GetForecast(ctx context.Context, lat, lon string) (*darksky.Forecast, error)
}

// Alexa handles requests made from the Amazon Echo
//
// TODO will it rain? "precipProbability": 0.73 && "precipType": "rain",
func Alexa(alexaAPI alexa.API, db *geo.DB, dsapi darkskyAPI) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) (int, interface{}) {
		var req alexa.Request
		if err := tea.Body(r, &req); err != nil {
			return tea.Error(400, errors.Wrap(err, "invalid request").Error())
		}

		loc := getLocation(r.Context(), &req, alexaAPI, db)

		forecast, err := dsapi.GetForecast(r.Context(), loc.Latitude, loc.Longitude)
		if err != nil {
			log.WithError(err).Error("failed to get forecast from darksky")
			return tea.StatusError(500)
		}

		question := alexa.ParseWeatherRequest(req.Request)
		response := speech.Speak(*forecast, question)

		return 200, alexa.ResponseText(response)
	}
	return tea.Handler(fn)
}

func getLocation(ctx context.Context, req *alexa.Request, api alexa.API, db *geo.DB) location {
	deviceID, accessToken := req.Context.System.Device.DeviceID, req.Context.System.APIAccessToken
	ll := log.WithFields(log.Fields{
		"device.id":      deviceID,
		"hasAccessToken": accessToken != "",
	})

	var zip string
	if deviceID == "" || accessToken == "" {
		ll.Info("no device info was found in request")
		return defaultLocation
	}
	var err error
	zip, err = api.DeviceZip(ctx, deviceID, accessToken)
	if err != nil {
		ll.WithError(err).Error("failed to retrieve zipcode for device, using default")
	}

	if lat, lon, ok := db.Lookup(zip); ok {
		return location{Latitude: lat, Longitude: lon}
	}
	ll.WithField("zip", zip).Error("failed to retrieve zipcode from geoDB. Using default location")
	return defaultLocation
}
