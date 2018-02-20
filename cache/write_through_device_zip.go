package cache

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/alexa"
)

// WriteThroughDevice is a cache layer that has a fallback layer
type WriteThroughDevice struct {
	cache Device
	api   alexa.API
}

// NewWriteThroughDevice creates a new WriteThrough cache
func NewWriteThroughDevice(cache Device, api alexa.API) *WriteThroughDevice {
	return &WriteThroughDevice{
		cache: cache,
		api:   api,
	}
}

// DeviceZip first tries to retrieve a cached result and falls back to
// directly fetching to the API. If the API is used then results are
// stored in the cache store
func (w *WriteThroughDevice) DeviceZip(ctx context.Context, accessToken, deviceID string) (string, error) {
	ll := log.WithFields(log.Fields{
		"component": "zipcache",
		"device.id": clip(deviceID, 25),
	})

	cached, err := w.cache.DeviceZip(ctx, deviceID)
	if err != nil {
		ll.WithError(err).Error("failed to get cached zip")
	}
	if cached != "" {
		ll.Info("cache hit")
		return cached, nil
	}
	ll.Info("cache miss")

	result, err := w.api.DeviceZip(ctx, deviceID, accessToken)
	if err != nil {
		return "", err
	}
	if result != "" {
		go func(ll log.Interface) {
			err := w.cache.PutDeviceZip(context.Background(), deviceID, result)
			if err != nil {
				ll.WithError(err).Error("failed to put cache")
			}
		}(ll)
	}

	return result, nil
}

func clip(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return fmt.Sprintf("%s...", s[:max])
}
