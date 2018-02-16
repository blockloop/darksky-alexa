package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/blockloop/darksky-alexa/darksky"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

// RedisPool is a pool of redis connections
type RedisPool interface {
	Get() redis.Conn
}

// Forecast is a layer which caches darksky forecast results
type Forecast interface {
	GetForecast(ctx context.Context, lat, lon string) (*darksky.Forecast, error)
	PutForecast(ctx context.Context, lat, lon string, f *darksky.Forecast) error
}

// Device is a device cache which caches deviceIDs to zipcodes to prevent the
// constant lookup of zipcode for deviceID
type Device interface {
	DeviceZip(ctx context.Context, deviceID string) (string, error)
	PutDeviceZip(ctx context.Context, deviceID, zip string) error
}

// NewRedis creates a new redis cache store with the default TTL of 15m
func NewRedis(pool RedisPool) *Redis {
	r := &Redis{
		pool: pool,
	}
	r.SetTTL(time.Minute * 15)
	return r
}

// Redis is a cache that uses redis
type Redis struct {
	pool RedisPool
	ttl  int64
}

// SetTTL sets the TTL used for every cache record in Redis
func (c *Redis) SetTTL(dur time.Duration) {
	c.ttl = int64(dur.Seconds())
}

// DeviceZip retrieves a zipcode for a deviceID
func (c *Redis) DeviceZip(ctx context.Context, deviceID string) (string, error) {
	con := c.pool.Get()
	defer con.Close()

	zip, err := redis.String(con.Do("HGET", "devices", deviceID))
	if err != nil {
		if err == redis.ErrNil {
			return "", nil
		}
		return "", errors.Wrap(err, "failed to get device")
	}

	return zip, nil
}

// PutDeviceZip stores a zipcode for a deviceID
func (c *Redis) PutDeviceZip(ctx context.Context, deviceID, zip string) error {
	con := c.pool.Get()
	defer con.Close()

	_, err := con.Do("HSETNX", "devices", deviceID, zip)
	return errors.Wrap(err, "failed to set cache")
}

// GetForecast retrieves a cache forecast from the redis store. If
// no cache exists then nil, nil is returned.
func (c *Redis) GetForecast(ctx context.Context, lat, lon string) (*darksky.Forecast, error) {
	con := c.pool.Get()
	defer con.Close()

	key := fmt.Sprintf("forecast:%s,%s", lat, lon)
	raw, err := redis.Bytes(con.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get forecast")
	}

	var forecast darksky.Forecast
	err = json.Unmarshal(raw, &forecast)
	return &forecast, errors.Wrap(err, "failed to parse JSON from redis")
}

// PutForecast stores a forecast to the redis store
func (c *Redis) PutForecast(ctx context.Context, lat, lon string, f *darksky.Forecast) error {
	con := c.pool.Get()
	defer con.Close()

	encoded, err := json.Marshal(f)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	key := fmt.Sprintf("forecast:%s,%s", lat, lon)
	con.Send("MULTI")
	con.Send("SET", key, encoded)
	con.Send("EXPIRE", key, c.ttl)
	_, err = con.Do("EXEC")

	return errors.Wrap(err, "failed to set cache")
}
