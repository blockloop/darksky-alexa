package speakers

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/blockloop/darksky-alexa/alexa"
	"github.com/blockloop/darksky-alexa/darksky"
	nowutil "github.com/jinzhu/now"
)

func TestPrecipitationPasses(t *testing.T) {
	t.Skip()
}

func TestPrecipitationTonight(t *testing.T) {
	fl, err := os.Open("json.json")
	require.NoError(t, err)

	var f darksky.Forecast
	require.NoError(t, json.NewDecoder(fl).Decode(&f))

	now := nowutil.New(time.Now())
	ts := &alexa.TimeSpan{
		End:   now.EndOfDay(),
		Start: now.Time,
	}
	p := Precipitation{}.dataPoints(ts, &f)
}
