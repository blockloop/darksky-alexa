package darksky

type Forecast RawForecast

type RawForecast struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Currently struct {
		UnixTime             int     `json:"time"`
		Summary              string  `json:"summary"`
		Icon                 string  `json:"icon"`
		NearestStormDistance int     `json:"nearestStormDistance"`
		PrecipIntensity      float64 `json:"precipIntensity"`
		PrecipIntensityError float64 `json:"precipIntensityError"`
		PrecipProbability    float64 `json:"precipProbability"`
		PrecipType           string  `json:"precipType"`
		Temperature          float64 `json:"temperature"`
		ApparentTemperature  float64 `json:"apparentTemperature"`
		DewPoint             float64 `json:"dewPoint"`
		Humidity             float64 `json:"humidity"`
		Pressure             float64 `json:"pressure"`
		WindSpeed            float64 `json:"windSpeed"`
		WindGust             float64 `json:"windGust"`
		WindBearing          int     `json:"windBearing"`
		CloudCover           float64 `json:"cloudCover"`
		UvIndex              int     `json:"uvIndex"`
		Visibility           float64 `json:"visibility"`
		Ozone                float64 `json:"ozone"`
	} `json:"currently"`
	Minutely struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
		Data    []struct {
			UnixTime             int     `json:"time"`
			PrecipIntensity      float64 `json:"precipIntensity"`
			PrecipIntensityError float64 `json:"precipIntensityError"`
			PrecipProbability    float64 `json:"precipProbability"`
			PrecipType           string  `json:"precipType"`
		} `json:"data"`
	} `json:"minutely"`
	Hourly struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
		Data    []struct {
			UnixTime            int     `json:"time"`
			Summary             string  `json:"summary"`
			Icon                string  `json:"icon"`
			PrecipIntensity     float64 `json:"precipIntensity"`
			PrecipProbability   float64 `json:"precipProbability"`
			PrecipType          string  `json:"precipType"`
			Temperature         float64 `json:"temperature"`
			ApparentTemperature float64 `json:"apparentTemperature"`
			DewPoint            float64 `json:"dewPoint"`
			Humidity            float64 `json:"humidity"`
			Pressure            float64 `json:"pressure"`
			WindSpeed           float64 `json:"windSpeed"`
			WindGust            float64 `json:"windGust"`
			WindBearing         int     `json:"windBearing"`
			CloudCover          float64 `json:"cloudCover"`
			UvIndex             int     `json:"uvIndex"`
			Visibility          float64 `json:"visibility"`
			Ozone               float64 `json:"ozone"`
		} `json:"data"`
	} `json:"hourly"`
	Daily struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
		Data    []struct {
			UnixTime                        int     `json:"time"`
			Summary                         string  `json:"summary"`
			Icon                            string  `json:"icon"`
			SunriseUnixTime                 int     `json:"sunriseTime"`
			SunsetUnixTime                  int     `json:"sunsetTime"`
			MoonPhase                       float64 `json:"moonPhase"`
			PrecipIntensity                 float64 `json:"precipIntensity"`
			PrecipIntensityMax              float64 `json:"precipIntensityMax"`
			PrecipIntensityMaxUnixTime      int     `json:"precipIntensityMaxTime"`
			PrecipProbability               float64 `json:"precipProbability"`
			PrecipType                      string  `json:"precipType"`
			TemperatureHigh                 float64 `json:"temperatureHigh"`
			TemperatureHighUnixTime         int     `json:"temperatureHighTime"`
			TemperatureLow                  float64 `json:"temperatureLow"`
			TemperatureLowUnixTime          int     `json:"temperatureLowTime"`
			ApparentTemperatureHigh         float64 `json:"apparentTemperatureHigh"`
			ApparentTemperatureHighUnixTime int     `json:"apparentTemperatureHighTime"`
			ApparentTemperatureLow          float64 `json:"apparentTemperatureLow"`
			ApparentTemperatureLowUnixTime  int     `json:"apparentTemperatureLowTime"`
			DewPoint                        float64 `json:"dewPoint"`
			Humidity                        float64 `json:"humidity"`
			Pressure                        float64 `json:"pressure"`
			WindSpeed                       float64 `json:"windSpeed"`
			WindGust                        float64 `json:"windGust"`
			WindGustUnixTime                int     `json:"windGustTime"`
			WindBearing                     int     `json:"windBearing"`
			CloudCover                      float64 `json:"cloudCover"`
			UvIndex                         int     `json:"uvIndex"`
			UvIndexUnixTime                 int     `json:"uvIndexTime"`
			Visibility                      float64 `json:"visibility"`
			Ozone                           float64 `json:"ozone"`
			TemperatureMin                  float64 `json:"temperatureMin"`
			TemperatureMinUnixTime          int     `json:"temperatureMinTime"`
			TemperatureMax                  float64 `json:"temperatureMax"`
			TemperatureMaxUnixTime          int     `json:"temperatureMaxTime"`
			ApparentTemperatureMin          float64 `json:"apparentTemperatureMin"`
			ApparentTemperatureMinUnixTime  int     `json:"apparentTemperatureMinTime"`
			ApparentTemperatureMax          float64 `json:"apparentTemperatureMax"`
			ApparentTemperatureMaxUnixTime  int     `json:"apparentTemperatureMaxTime"`
		} `json:"data"`
	} `json:"daily"`
	Alerts []struct {
		Title       string `json:"title"`
		UnixTime    int    `json:"time"`
		Expires     int    `json:"expires"`
		Description string `json:"description"`
		URI         string `json:"uri"`
	} `json:"alerts"`
}
