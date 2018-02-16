package alexa

import (
	"fmt"
	"time"
)

// TimeSpan is a period of time indicating a start time,
// a Span, and a count
type TimeSpan struct {
	Start time.Time
	End   time.Time
}

func (ts TimeSpan) String() string {
	return fmt.Sprintf("Start: %s, End: %s, Dur: %s",
		ts.Start.Format(time.Stamp),
		ts.End.Format(time.Stamp),
		ts.Diff())
}

func (ts TimeSpan) Diff() time.Duration {
	return ts.End.Sub(ts.Start)
}

// IsZero returns true if TimeSpan is zero value
func (ts TimeSpan) IsZero() bool {
	return ts.Start.IsZero()
}
