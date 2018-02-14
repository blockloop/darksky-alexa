package alexa

import "time"

// TimeSpan is a period of time indicating a start time,
// a Span, and a count
type TimeSpan struct {
	Start time.Time
	Span  time.Duration
}

// IsZero returns true if TimeSpan is zero value
func (ts TimeSpan) IsZero() bool {
	return ts.Start.IsZero()
}
