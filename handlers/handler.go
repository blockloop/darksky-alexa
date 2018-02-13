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
func Alexa(alexaAPI alexa.API, geoDB *geo.DB, dsapi darkskyAPI) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) (int, interface{}) {
		var req alexa.Request
		if err := tea.Body(r, &req); err != nil {
			return tea.Error(400, errors.Wrap(err, "invalid request").Error())
		}

		zip, err := alexaAPI.DeviceZip(r.Context(), req.Context.System.Device.DeviceID, req.Context.System.APIAccessToken)
		if err != nil {
			log.WithError(err).Error("failed to retrieve zipcode for device")
			return tea.StatusError(500)
		}
		lat, lon, ok := geoDB.Lookup(zip)
		if !ok {
			log.Error("failed to retrieve zipcode from geoDB. Using default location")
			lat, lon = defaultLat, defaultLon
		}

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