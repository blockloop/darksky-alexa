package handlers

// defaultLocation is used when we cannot determine the location
// of the device
var defaultLocation = location{
	Latitude:  "32.857112",
	Longitude: "-96.431210",
	Zip:       "75032",
}

type location struct {
	Zip       string
	Latitude  string
	Longitude string
}
