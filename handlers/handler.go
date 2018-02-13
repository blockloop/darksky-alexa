package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/darksky-alexa/geo"
	"github.com/blockloop/tea"
	"github.com/pkg/errors"
)

const (
	// default location is set to zipcode 75032 because that's where I live
	defaultLat = "32.857112"
	defaultLon = "-96.431210"
)

// Ping is a simple handler that responds for healthchecks
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"pong": ""}`))
}

type darkskyAPI interface {
	GetForecast(ctx context.Context, lat, lon string) (*darksky.Forecast, error)
}

// Alexa handles requests made from the Amazon Echo
func Alexa(alexaAPI alexa.API, db *geo.DB, dsapi darkskyAPI) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) (int, interface{}) {
		var req alexa.Request
		if err := tea.Body(r, &req); err != nil {
			return tea.Error(400, errors.Wrap(err, "invalid request").Error())
		}

		lat, lon := getLocation(r.Context(), &req, alexaAPI, db)

		return getWeather(r.Context(), lat, lon, dsapi)
	}
	return tea.Handler(fn)
}

func getWeather(ctx context.Context, lat, lon string, api darkskyAPI) (int, interface{}) {
	forecast, err := api.GetForecast(ctx, lat, lon)
	if err != nil {
		log.WithError(err).Error("failed to get forecast from darksky")
		return tea.StatusError(500)
	}

	msg := fmt.Sprintf("currently %.0f° and %s with a high of %.0f° and a low of %.0f°",
		forecast.Hourly.Data[0].Temperature,
		forecast.Currently.Summary,
		forecast.Daily.Data[0].TemperatureHigh,
		forecast.Daily.Data[0].TemperatureLow,
	)
	return 200, alexa.ResponseText(msg)
}

func getLocation(ctx context.Context, req *alexa.Request, api alexa.API, db *geo.DB) (lat, long string) {
	deviceID, apiHost, accessToken :=
		req.Context.System.Device.DeviceID, req.Context.System.APIEndpoint, req.Context.System.APIAccessToken
	ll := log.WithFields(log.Fields{
		"device.id":      deviceID,
		"api.endpoint":   apiHost,
		"hasAccessToken": accessToken != "",
	})

	var zip string
	if deviceID == "" || accessToken == "" {
		ll.Info("no device info was found in request")
		return defaultLat, defaultLon
	}
	var err error
	zip, err = api.DeviceZip(ctx, deviceID, apiHost, accessToken)
	if err != nil {
		ll.WithError(err).Error("failed to retrieve zipcode for device, using default")
	}

	if lat, lon, ok := db.Lookup(zip); ok {
		return lat, lon
	}
	ll.WithField("zip", zip).Error("failed to retrieve zipcode from geoDB. Using default location")
	return defaultLat, defaultLon
}
