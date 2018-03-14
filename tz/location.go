package tz

import "time"

type Location struct {
	Latitude  string
	Longitude string
	Timezone  *time.Location
	Zipcode   string
}
