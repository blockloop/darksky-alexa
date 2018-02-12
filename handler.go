package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/tea"
)

func handler(cache Cache, api *darksky.API) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) (int, interface{}) {
		// var req Request
		// if err := tea.Body(r, &req); err != nil {
		// 	return tea.Error(400, errors.Wrap(err, "invalid request").Error())
		// }
		// return receive(r.Context(), cache, api, req.Values())
		return receive(r.Context(), cache, api)
	}
	return tea.Handler(fn)
}

// func receive(ctx context.Context, cache Cache, api *darksky.API, req RequestValues) (int, interface{}) {
func receive(ctx context.Context, cache Cache, api *darksky.API) (int, interface{}) {
	forecast, err := api.GetForecast(ctx)
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
	return 200, ResponseText(msg)
}
