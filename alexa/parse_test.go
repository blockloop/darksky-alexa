package alexa

import (
	"testing"
	"time"

	"github.com/jinzhu/now"
	"github.com/stretchr/testify/assert"
)

func TestParseWeek(t *testing.T) {
	table := map[string]time.Time{
		"2018-W1":     time.Date(2018, time.January, 1, 0, 0, 0, 0, time.Local),
		"2018-W10":    time.Date(2018, time.March, 5, 0, 0, 0, 0, time.Local),
		"2018-W1-WE":  time.Date(2018, time.January, 6, 0, 0, 0, 0, time.Local),
		"2018-W10-WE": time.Date(2018, time.March, 10, 0, 0, 0, 0, time.Local),
	}

	for k, exp := range table {
		ts := parseWeek(k)
		assert.Equal(t, exp.Format(yyyymmdd), ts.Start.Format(yyyymmdd), k)
	}
}

func TestParseDayDefaultsToNowWhenNotSpecific(t *testing.T) {
	tests := []string{"2018", "201X", "WI", "SP", "SU", "FA", ""}

	today := now.New(time.Now())
	exp := TimeSpan{
		End:   today.EndOfDay(),
		Start: today.BeginningOfDay(),
	}
	for _, test := range tests {
		ts := parseDay(test)
		assert.Equal(t, exp.Start.String(), ts.Start.String())
		assert.Equal(t, exp.End.String(), ts.End.String())
	}
}

func TestParseDayUsesWeekSpanWhenSpecifiesMonth(t *testing.T) {
	test := "2018-01"
	ts := parseDay(test)
	eow := now.New(ts.Start).EndOfWeek()
	assert.Equal(t, eow.String(), ts.End.String(), "Span mismatch: %q", test)
}
