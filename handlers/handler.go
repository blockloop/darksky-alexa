package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/darksky-alexa/geo"
	"github.com/blockloop/darksky-alexa/pollen"
	"github.com/blockloop/darksky-alexa/speech"
	"github.com/blockloop/darksky-alexa/tz"
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

type pollenAPI interface {
	GetPollen(ctx context.Context, loc *tz.Location) (*pollen.Forecast, error)
}

// Alexa handles requests made from the Amazon Echo
func Alexa(alexaAPI alexa.API, db *geo.DB, dsapi darkskyAPI, papi pollenAPI) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) (int, interface{}) {
		var req alexa.Request
		if err := tea.Body(r, &req); err != nil {
			log.WithError(err).Info("received bad request")
			if err == io.EOF {
				err = errors.New("missing request body")
			}
			return tea.Error(400, errors.Wrap(err, "invalid request").Error())
		}
		ll := log.WithFields(log.Fields{
			"condition": req.Request.Intent.Slots.Condition.Value,
			"day":       req.Request.Intent.Slots.Day.Value,
			"time":      req.Request.Intent.Slots.Time.Value,
		})
		ll.Info("request")

		loc := getLocation(r.Context(), &req, alexaAPI, db)

		forecast, err := dsapi.GetForecast(r.Context(), loc.Latitude, loc.Longitude)
		if err != nil {
			log.WithError(err).Error("failed to get forecast from darksky")
			return tea.StatusError(500)
		}

		pol, err := papi.GetPollen(r.Context(), loc)
		if err != nil {
			log.WithError(err).Error("failed to get pollen data")
			return tea.StatusError(500)
		}

		q := alexa.ParseWeatherRequest(req.Request, loc.Timezone)
		ll.WithField("parsed", q.String()).Info("parsed request")
		response := speech.Speak(loc, forecast, pol, q)

		return 200, alexa.ResponseText(response)
	}
	return tea.Handler(fn)
}

func getLocation(ctx context.Context, req *alexa.Request, api alexa.API, db *geo.DB) *tz.Location {
	deviceID, accessToken := req.Context.System.Device.DeviceID, req.Context.System.APIAccessToken
	ll := log.WithFields(log.Fields{
		"device.id":   clip(deviceID, 25),
		"accessToken": clip(accessToken, 25),
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

	loc, err := db.LookupZip(ctx, zip)
	if err != nil {
		ll.WithField("zip", zip).Error("failed to retrieve zipcode from geoDB. Using default location")
		return defaultLocation
	}
	return loc
}

func clip(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return fmt.Sprintf("%s...", s[:max])
}
