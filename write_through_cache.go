package main

import (
	"context"

	"github.com/blockloop/darksky-alexa/darksky"
)

type ForecastAPI interface {
	GetForecast(ctx context.Context) (*darksky.Forecast, error)
}

var _ ForecastAPI = (*WriteThroughCache)(nil)

type WriteThroughCache struct {
	cache Cache
	api   *darksky.API
}

func (w *WriteThroughCache) GetForecast(ctx context.Context) (*darksky.Forecast, error) {
	return w.api.GetForecast(ctx)
}
