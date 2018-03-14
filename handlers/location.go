package handlers

import (
	"time"

	"github.com/apex/log"
	"github.com/blockloop/darksky-alexa/tz"
)

// defaultLocation is used when we cannot determine the location
// of the device
var defaultLocation = &tz.Location{
	Latitude:  "32.857112",
	Longitude: "-96.431210",
	Zipcode:   "75032",
	Timezone:  mustTimezone("America/Chicago"),
}

func mustTimezone(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		log.WithField("timezone", name).
			WithError(err).
			Fatal("failed to load timezone")
	}
	return loc
}
