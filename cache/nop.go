package cache

import (
	"context"

	"github.com/blockloop/darksky-alexa/darksky"
)

// Nop is a cache that does nothing
type Nop struct {
}

func (c *Nop) GetForecast(ctx context.Context, lon, lat string) (*darksky.Forecast, error) {
	return nil, nil
}

func (c *Nop) PutForecast(ctx context.Context, lon, lat string, f *darksky.Forecast) error {
	return nil
}
