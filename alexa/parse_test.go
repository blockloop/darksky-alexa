package alexa

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseWeek(t *testing.T) {
	format := "2006-01-02"
	table := map[string]time.Time{
		"2018-W1":     time.Date(2018, time.January, 1, 0, 0, 0, 0, time.Local),
		"2018-W10":    time.Date(2018, time.March, 5, 0, 0, 0, 0, time.Local),
		"2018-W1-WE":  time.Date(2018, time.January, 6, 0, 0, 0, 0, time.Local),
		"2018-W10-WE": time.Date(2018, time.March, 10, 0, 0, 0, 0, time.Local),
	}

	for k, exp := range table {
		start := parseWeek(k)
		assert.Equal(t, exp.Format(format), start.Format(format), k)
	}
}

func TestParseDayDefaultsToNowWhenNotSpecific(t *testing.T) {
	tests := []string{"2018", "201X", "WI", "SP", "SU", "FA"}

	exp := TimeSpan{
		Count: 0,
		Span:  Day,
		Start: time.Now().UTC(),
	}
	for _, test := range tests {
		ts := parseDay(test)
		assert.Equal(t, exp.Start.Format("2006-01-02"), ts.Start.Format("2006-01-02"))
		assert.Equal(t, exp.Count, ts.Count, "Count mismatch: %q", test)
		assert.Equal(t, exp.Span.String(), ts.Span.String(), "Span mismatch: %q", test)
	}
}

func TestParseDayUsesWeekSpanWhenSpecifiesMonth(t *testing.T) {
	test := "2018-01"
	ts := parseDay(test)
	assert.Equal(t, Week.String(), ts.Span.String(), "Span mismatch: %q", test)
}
