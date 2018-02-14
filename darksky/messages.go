package darksky

type Forecast RawForecast

type RawForecast struct {
	Alerts    []Alert    `json:"alerts"`
	Currently DataPoint  `json:"currently"`
	Daily     DataPoints `json:"daily"`
	Hourly    DataPoints `json:"hourly"`
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
	Minutely  DataPoints `json:"minutely"`
	Timezone  string     `json:"timezone"`
}

type DataPoints struct {
	Summary string      `json:"summary"`
	Icon    string      `json:"icon"`
	Data    []DataPoint `json:"data"`
}

type Currently struct {
}

type DataPoint struct {
	ApparentTemperature             float64 `json:"apparentTemperature"`
	ApparentTemperatureHigh         float64 `json:"apparentTemperatureHigh"`
	ApparentTemperatureHighUnixTime int64   `json:"apparentTemperatureHighTime"`
	ApparentTemperatureLow          float64 `json:"apparentTemperatureLow"`
	ApparentTemperatureLowUnixTime  int64   `json:"apparentTemperatureLowTime"`
	ApparentTemperatureMax          float64 `json:"apparentTemperatureMax"`
	ApparentTemperatureMaxUnixTime  int64   `json:"apparentTemperatureMaxTime"`
	ApparentTemperatureMin          float64 `json:"apparentTemperatureMin"`
	ApparentTemperatureMinUnixTime  int64   `json:"apparentTemperatureMinTime"`
	CloudCover                      float64 `json:"cloudCover"`
	DewPoint                        float64 `json:"dewPoint"`
	Humidity                        float64 `json:"humidity"`
	Icon                            string  `json:"icon"`
	MoonPhase                       float64 `json:"moonPhase"`
	NearestStormDistance            int     `json:"nearestStormDistance"`
	Ozone                           float64 `json:"ozone"`
	PrecipIntensity                 float64 `json:"precipIntensity"`
	PrecipIntensityError            float64 `json:"precipIntensityError"`
	PrecipIntensityMax              float64 `json:"precipIntensityMax"`
	PrecipIntensityMaxUnixTime      int64   `json:"precipIntensityMaxTime"`
	PrecipProbability               float64 `json:"precipProbability"`
	PrecipType                      string  `json:"precipType"`
	Pressure                        float64 `json:"pressure"`
	Summary                         string  `json:"summary"`
	SunriseUnixTime                 int64   `json:"sunriseTime"`
	SunsetUnixTime                  int64   `json:"sunsetTime"`
	Temperature                     float64 `json:"temperature"`
	TemperatureHigh                 float64 `json:"temperatureHigh"`
	TemperatureHighUnixTime         int64   `json:"temperatureHighTime"`
	TemperatureLow                  float64 `json:"temperatureLow"`
	TemperatureLowUnixTime          int64   `json:"temperatureLowTime"`
	TemperatureMax                  float64 `json:"temperatureMax"`
	TemperatureMaxUnixTime          int64   `json:"temperatureMaxTime"`
	TemperatureMin                  float64 `json:"temperatureMin"`
	TemperatureMinUnixTime          int64   `json:"temperatureMinTime"`
	UnixTime                        int64   `json:"time"`
	UvIndex                         int     `json:"uvIndex"`
	UvIndexUnixTime                 int64   `json:"uvIndexTime"`
	Visibility                      float64 `json:"visibility"`
	WindBearing                     int     `json:"windBearing"`
	WindGust                        float64 `json:"windGust"`
	WindGustUnixTime                int64   `json:"windGustTime"`
	WindSpeed                       float64 `json:"windSpeed"`
}

type Alert struct {
	Description string `json:"description"`
	Expires     int64  `json:"expires"`
	Title       string `json:"title"`
	URI         string `json:"uri"`
	UnixTime    int64  `json:"time"`
}
