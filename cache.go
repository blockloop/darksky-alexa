package main

import (
	"context"
	"time"

	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/garyburd/redigo/redis"
)

type Cache interface {
	GetForecast(ctx context.Context, day time.Time, city string) (*darksky.Forecast, error)
	PutForecast(ctx context.Context, day time.Time, city string, f *darksky.Forecast) error
}

// RedisCache is a cache that uses redis
type RedisCache struct {
	pool *redis.Pool
}

func (c *RedisCache) GetForecast(ctx context.Context, day time.Time, city string) (*darksky.Forecast, error) {
	return nil, nil
}

func (c *RedisCache) PutForecast(ctx context.Context, day time.Time, city string, f *darksky.Forecast) error {
	return nil
}

// NopCache is a cache that does nothing
type NopCache struct {
}

func (*NopCache) GetForecast(ctx context.Context, day time.Time, city string) (*darksky.Forecast, error) {
	return nil, nil
}

func (*NopCache) PutForecast(ctx context.Context, day time.Time, city string, f *darksky.Forecast) error {
	return nil
}
