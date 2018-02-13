package cache

import (
	"context"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/pkg/errors"
)

// WriteThroughForecast is a cache layer that has a fallback layer
type WriteThroughForecast struct {
	cache Forecast
	api   *darksky.API
}

// NewWriteThroughForecast creates a new WriteThrough cache
func NewWriteThroughForecast(cache Forecast, api *darksky.API) *WriteThroughForecast {
	return &WriteThroughForecast{
		cache: cache,
		api:   api,
	}
}

// GetForecast first tries to retrieve a cached result and falls back to
// directly fetching to the API. If the API is used then results are
// stored in the cache store
func (w *WriteThroughForecast) GetForecast(ctx context.Context, lat, lon string) (*darksky.Forecast, error) {
	ll := log.WithFields(log.Fields{
		"component": "forecastcache",
		"latitude":  lat,
		"longitude": lon,
	})

	forecast, err := w.cache.GetForecast(ctx, lat, lon)
	if err != nil {
		ll.WithError(err).Error("failed to get cached forecast")
	}
	if forecast != nil {
		ll.Info("cache hit")
		return forecast, nil
	}
	ll.Info("cache miss")

	forecast, err = w.api.GetForecast(ctx, lat, lon)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch forecast from API")
	}
	if forecast != nil {
		go func(ll log.Interface) {
			err := w.cache.PutForecast(context.Background(), lat, lon, forecast)
			if err != nil {
				ll.WithError(err).Error("failed to put cache")
			}
		}(ll)
	}

	return forecast, nil
}
