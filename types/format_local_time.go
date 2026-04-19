package types

import (
	"fmt"
	"time"
)

// localTimeFormats lists the parser layouts accepted by parseLocalTime, in
// the order they are tried. The ".999999999" suffix makes fractional seconds
// optional on parse and trims trailing zeros on format.
var localTimeFormats = []string{
	"15:04:05.999999999",
	"15:04",
}

// formatLocalTime renders a wall-clock time-of-day as "15:04:05[.fff]" with
// no timezone designator. Fractional seconds are emitted only when non-zero;
// trailing zeros are trimmed.
func formatLocalTime(lt LocalTime) string {
	return time.Date(0, 1, 1,
		lt.Hour, lt.Minute, lt.Second, lt.Nanosec,
		time.UTC,
	).Format("15:04:05.999999999")
}

// parseLocalTime parses a wall-clock time-of-day. Inputs carrying a timezone
// designator are rejected with an error.
func parseLocalTime(s string) (LocalTime, error) {
	if hasOffsetSuffix(s) {
		return LocalTime{}, fmt.Errorf("LocalTime must not carry timezone info: %q", s)
	}
	for _, f := range localTimeFormats {
		parsed, err := time.Parse(f, s)
		if err == nil {
			return localTimeFromTime(parsed), nil
		}
	}
	return LocalTime{}, fmt.Errorf("cannot parse LocalTime from %q", s)
}

// localTimeFromTime extracts the time-of-day components of t and discards
// its date and location.
func localTimeFromTime(t time.Time) LocalTime {
	return LocalTime{
		Hour:    t.Hour(),
		Minute:  t.Minute(),
		Second:  t.Second(),
		Nanosec: t.Nanosecond(),
	}
}
