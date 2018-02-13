package handlers

// defaultLocation is used when we cannot determine the location
// of the device
var defaultLocation = location{
	Latitude:  "32.857112",
	Longitude: "-96.431210",
}

type location struct {
	Latitude  string
	Longitude string
}
