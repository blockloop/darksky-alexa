package ratelimiter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/cache"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

// New creates a new rate limiter
func New(pool cache.RedisPool, totalPerDay, ipPerDay int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			today := now.Format("2006-01-02")
			totalKey := fmt.Sprintf("%s.%s:requests", today, r.RemoteAddr)
			ipKey := fmt.Sprintf("%s:requests", today)

			eod := eodSeconds(now)

			totalCount, err := incrLimit(pool.Get(), totalKey, eod)
			if err != nil {
				log.WithError(err).WithField("key", totalKey).Error("failed retrieving count")
			}

			ipCount, err := incrLimit(pool.Get(), ipKey, eod)
			if err != nil {
				log.WithError(err).WithField("key", ipKey).Error("failed retrieving count")
			}

			if totalCount >= totalPerDay {
				log.WithFields(log.Fields{
					"key":   totalKey,
					"count": totalCount,
					"limit": totalPerDay,
				}).Warn("too many requests")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			if ipCount >= ipPerDay {
				log.WithFields(log.Fields{
					"key":   ipKey,
					"count": ipCount,
					"limit": ipPerDay,
				}).Warn("too many requests")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func incrLimit(con redis.Conn, key string, ttl int) (int64, error) {
	defer con.Close()
	con.Send("MULTI")
	con.Send("INCR", key)
	con.Send("EXPIRE", key, ttl)
	vals, err := redis.Int64s(con.Do("EXEC"))
	if err != nil {
		return 0, errors.Wrap(err, "failed to INCR")
	}
	return vals[0], nil
}

func eodSeconds(now time.Time) int {
	year, month, day := now.Date()
	eod := time.Date(year, month, day, 0, 0, 0, 0, now.Location()).Add(time.Hour * 24)
	return int(eod.Sub(now).Seconds())
}
