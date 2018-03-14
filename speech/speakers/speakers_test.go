package speakers

import (
	"fmt"
	"testing"
	"time"

	timeclock "github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHumanDayTable(t *testing.T) {
	// Sunday
	mockNow := time.Date(2018, time.February, 12, 1, 0, 0, 0, time.Local)
	mockclock := timeclock.NewMock()
	mockclock.Set(mockNow)
	clock = mockclock

	table := map[string]time.Time{
		"today":             mockNow,
		"tomorrow":          mockNow.Add(time.Hour * 24),
		"Sunday":            mockNow.AddDate(0, 0, 6).Add(time.Hour),
		"Wednesday, Feb 21": mockNow.AddDate(0, 0, 9),
	}

	i := 0
	for expected, test := range table {
		assert.Equal(t, expected, humanDay(test), fmt.Sprintf("test #%d", i))
		i++
	}
}

func TestSameDayReturnsTrueWhenDifferentTimezone(t *testing.T) {
	cdt, err := time.LoadLocation("America/Chicago")
	require.NoError(t, err)

	thetime, err := time.ParseInLocation(time.RFC3339, "2018-03-13T22:00:00-05:00", cdt)
	require.NoError(t, err)

	utc := thetime.UTC()

	assert.Truef(t, sameDay(thetime, utc), "%s == %s", thetime, utc)
}
