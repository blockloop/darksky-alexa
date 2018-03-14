package cache

import (
	"context"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/pollen"
	"github.com/blockloop/darksky-alexa/tz"
	"github.com/pkg/errors"
)

// WriteThroughPollen is a cache layer that has a fallback layer
type WriteThroughPollen struct {
	cache Pollen
	api   *pollen.API
}

// NewWriteThroughPollen creates a new WriteThrough cache
func NewWriteThroughPollen(cache Pollen, api *pollen.API) *WriteThroughPollen {
	return &WriteThroughPollen{
		cache: cache,
		api:   api,
	}
}

// GetPollen first tries to retrieve a cached result and falls back to
// directly fetching to the API. If the API is used then results are
// stored in the cache store
func (w *WriteThroughPollen) GetPollen(ctx context.Context, loc *tz.Location) (*pollen.Forecast, error) {
	ll := log.WithFields(log.Fields{
		"component": "pollencache",
		"zip":       loc.Zipcode,
	})

	forecast, err := w.cache.GetPollen(ctx, loc.Zipcode)
	if err != nil {
		ll.WithError(err).Error("failed to get cached pollen")
	}
	if forecast != nil {
		ll.Info("cache hit")
		return forecast, nil
	}
	ll.Info("cache miss")

	forecast, err = w.api.GetPollen(ctx, loc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch pollen from API")
	}
	if forecast != nil {
		go func(ll log.Interface) {
			err := w.cache.PutPollen(context.Background(), loc.Zipcode, forecast)
			if err != nil {
				ll.WithError(err).Error("failed to put cache")
			}
		}(ll)
	}

	return forecast, nil
}
