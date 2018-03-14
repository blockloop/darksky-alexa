package geo

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/blockloop/darksky-alexa/tz"
	"github.com/pkg/errors"
)

var (
	selectgeo = sq.
		Select("latitude", "longitude", "zip", "timezone").
		From("geo")
)

// New creates a new geo DB that uses the given redis pool and
// triggers db.Load in a separate goroutine. Check DB.Loaded()
// to determine when the datastore is fully loaded
func New(db *sql.DB) *DB {
	return &DB{
		db: db,
	}
}

// DB stores geo info to a redis store
type DB struct {
	db *sql.DB
}

// LookupZip gets a latitude and longitude for a zipcode. If the result
// is found then ok is TRUE along with the results If no result is
// found then ok is FALSE
func (db *DB) LookupZip(ctx context.Context, zip string) (*tz.Location, error) {
	loc := new(tz.Location)
	var zone string

	err := selectgeo.
		Where(sq.Eq{"zip": zip}).
		Limit(1).
		RunWith(db.db).
		QueryRowContext(ctx).
		Scan(&loc.Latitude, &loc.Longitude, &loc.Zipcode, &zone)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query timezone")
	}

	loc.Timezone, err = time.LoadLocation(zone)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load location %q", zone)
	}

	return loc, nil
}
