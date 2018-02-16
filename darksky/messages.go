package darksky

import (
	"time"
)

// UnixTime is a timestamp that unmarshals from a unix timestamp
type UnixTime int64

func (t UnixTime) Time() time.Time {
	return time.Unix(int64(t), 0)
}

type Forecast struct {
	Alerts    []Alert   `json:"alerts,omitempty"`
	Currently DataPoint `json:"currently,omitempty"`
	Daily     DataBlock `json:"daily,omitempty"`
	Hourly    DataBlock `json:"hourly,omitempty"`
	Latitude  float64   `json:"latitude,omitempty"`
	Longitude float64   `json:"longitude,omitempty"`
	Minutely  DataBlock `json:"minutely,omitempty"`
	Timezone  string    `json:"timezone,omitempty"`
}

type DataBlock struct {
	Summary string      `json:"summary,omitempty"`
	Icon    string      `json:"icon,omitempty"`
	Data    []DataPoint `json:"data,omitempty"`
}

type DataPoint struct {
	ApparentTemperature         float64  `json:"apparentTemperature,omitempty"`
	ApparentTemperatureHigh     float64  `json:"apparentTemperatureHigh,omitempty"`
	ApparentTemperatureHighTime UnixTime `json:"apparentTemperatureHighTime,omitempty"`
	ApparentTemperatureLow      float64  `json:"apparentTemperatureLow,omitempty"`
	ApparentTemperatureLowTime  UnixTime `json:"apparentTemperatureLowTime,omitempty"`
	ApparentTemperatureMax      float64  `json:"apparentTemperatureMax,omitempty"`
	ApparentTemperatureMaxTime  UnixTime `json:"apparentTemperatureMaxTime,omitempty"`
	ApparentTemperatureMin      float64  `json:"apparentTemperatureMin,omitempty"`
	ApparentTemperatureMinTime  UnixTime `json:"apparentTemperatureMinTime,omitempty"`
	CloudCover                  float64  `json:"cloudCover,omitempty"`
	DewPoint                    float64  `json:"dewPoint,omitempty"`
	Humidity                    float64  `json:"humidity,omitempty"`
	Icon                        string   `json:"icon,omitempty"`
	MoonPhase                   float64  `json:"moonPhase,omitempty"`
	NearestStormDistance        int      `json:"nearestStormDistance,omitempty"`
	Ozone                       float64  `json:"ozone,omitempty"`
	PrecipIntensity             float64  `json:"precipIntensity,omitempty"`
	PrecipIntensityError        float64  `json:"precipIntensityError,omitempty"`
	PrecipIntensityMax          float64  `json:"precipIntensityMax,omitempty"`
	PrecipIntensityMaxTime      UnixTime `json:"precipIntensityMaxTime,omitempty"`
	PrecipProbability           float64  `json:"precipProbability,omitempty"`
	PrecipType                  string   `json:"precipType,omitempty"`
	Pressure                    float64  `json:"pressure,omitempty"`
	Summary                     string   `json:"summary,omitempty"`
	SunriseTime                 UnixTime `json:"sunriseTime,omitempty"`
	SunsetTime                  UnixTime `json:"sunsetTime,omitempty"`
	Temperature                 float64  `json:"temperature,omitempty"`
	TemperatureHigh             float64  `json:"temperatureHigh,omitempty"`
	TemperatureHighTime         UnixTime `json:"temperatureHighTime,omitempty"`
	TemperatureLow              float64  `json:"temperatureLow,omitempty"`
	TemperatureLowTime          UnixTime `json:"temperatureLowTime,omitempty"`
	TemperatureMax              float64  `json:"temperatureMax,omitempty"`
	TemperatureMaxTime          UnixTime `json:"temperatureMaxTime,omitempty"`
	TemperatureMin              float64  `json:"temperatureMin,omitempty"`
	TemperatureMinTime          UnixTime `json:"temperatureMinTime,omitempty"`
	Time                        UnixTime `json:"time,omitempty"`
	UvIndex                     int      `json:"uvIndex,omitempty"`
	UvIndexTime                 UnixTime `json:"uvIndexTime,omitempty"`
	Visibility                  float64  `json:"visibility,omitempty"`
	WindBearing                 int      `json:"windBearing,omitempty"`
	WindGust                    float64  `json:"windGust,omitempty"`
	WindGustTime                UnixTime `json:"windGustTime,omitempty"`
	WindSpeed                   float64  `json:"windSpeed,omitempty"`
}

type Alert struct {
	Description string   `json:"description,omitempty"`
	Expires     UnixTime `json:"expires,omitempty"`
	Title       string   `json:"title,omitempty"`
	URI         string   `json:"uri,omitempty"`
	Time        UnixTime `json:"time,omitempty"`
}
