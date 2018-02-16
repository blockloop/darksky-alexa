package geo

import (
	"bufio"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/cache"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

// New creates a new geo DB that uses the given redis pool and
// triggers db.Load in a separate goroutine. Check DB.Loaded()
// to determine when the datastore is fully loaded
func New(pool cache.RedisPool) *DB {
	db := NewUnloaded(pool)
	go db.load(zipcodesCSV)
	return db
}

// NewUnloaded creates a new unloaded geo DB. You should
func NewUnloaded(pool cache.RedisPool) *DB {
	return &DB{
		pool: pool,
		key:  "geo",
	}
}

// DB stores geo info to a redis store
type DB struct {
	key    string
	pool   cache.RedisPool
	loaded int32
}

// Loaded returns true if the datastore has been loaded
func (db *DB) Loaded() bool {
	return atomic.LoadInt32(&db.loaded) == 1
}

// SetKey sets the redis key to use for storage
func (db *DB) SetKey(key string) {
	db.key = key
}

// Lookup gets a latitude and longitude for a zipcode. If the result
// is found then ok is TRUE along with the results If no result is
// found then ok is FALSE
func (db *DB) Lookup(zipcode string) (lat, lon string, ok bool) {
	con := db.pool.Get()
	defer con.Close()

	raw, err := redis.String(con.Do("HGET", db.key, zipcode))
	if err != nil {
		if err == redis.ErrNil {
			return
		}
	}
	if raw == "" {
		return
	}
	splits := strings.Split(raw, ",")
	if len(splits) != 2 {
		log.WithFields(log.Fields{
			"zipcode": zipcode,
			"result":  raw,
		}).Error("bad result for zipcode stored in redis")
		return
	}

	ok = true
	lat, lon = splits[0], splits[1]
	return
}

func (db *DB) store(zipcode, lat, lon string) error {
	con := db.pool.Get()
	defer con.Close()

	val := fmt.Sprintf("%s,%s", lat, lon)
	_, err := con.Do("HSETNX", db.key, zipcode, val)
	// log.WithError(err).WithFields(log.Fields{
	// 	"key":       db.key,
	// 	"zip":       zipcode,
	// 	"lat":       lat,
	// 	"lon":       lon,
	// 	"composite": val,
	// }).Fatalf("setting key")
	return errors.Wrap(err, "failed to set hash key")
}

func (db *DB) load(csv string) {
	scanner := bufio.NewScanner(strings.NewReader(csv))
	ll := log.WithFields(log.Fields{
		"component": "geo loader",
	})

	start := time.Now()
	lineno := 0
	for scanner.Scan() {
		line := scanner.Text()
		ll := ll.WithFields(log.Fields{
			"line.text":   line,
			"line.number": lineno,
		})
		splits := strings.Split(line, ",")
		if len(splits) != 3 {
			ll.Error("line does not contain three items")
			continue
		}
		zip, lat, lon := splits[0], splits[1], splits[2]
		if len(zip) == 0 {
			ll.Error("zip is empty")
			continue
		}
		if len(lat) == 0 {
			ll.Error("lat is empty")
			continue
		}
		if len(lon) == 0 {
			ll.Error("lon is empty")
			continue
		}

		if err := db.store(zip, lat, lon); err != nil {
			ll.WithError(err).Error("failed to store line in redis")
			continue
		}

		lineno++
	}

	ll.WithField("duration", time.Since(start)).Info("datastore loaded")
	atomic.StoreInt32(&db.loaded, 1)
}
