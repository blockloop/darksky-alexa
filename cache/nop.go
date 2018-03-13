package cache

import (
	"context"

	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/blockloop/darksky-alexa/pollen"
)

// Nop is a cache that does nothing
type Nop struct {
}

func (n *Nop) GetForecast(ctx context.Context, lat string, lon string) (*darksky.Forecast, error) {
	return nil, nil
}

func (n *Nop) PutForecast(ctx context.Context, lat string, lon string, f *darksky.Forecast) error {
	return nil
}

func (n *Nop) DeviceZip(ctx context.Context, deviceID string) (string, error) {
	return "", nil
}

func (n *Nop) PutDeviceZip(ctx context.Context, deviceID string, zip string) error {
	return nil
}

func (n *Nop) GetPollen(ctx context.Context, zip string) (*pollen.Forecast, error) {
	return nil, nil
}

func (n *Nop) PutPollen(ctx context.Context, zip string, p *pollen.Forecast) error {
	return nil
}
