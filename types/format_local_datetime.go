package types

import (
	"fmt"
	"time"
)

// localDateTimeFormats lists the parser layouts accepted by parseLocalDateTime,
// in the order they are tried. The ".999999999" suffix makes fractional
// seconds optional on parse and trims trailing zeros on format.
var localDateTimeFormats = []string{
	"2006-01-02T15:04:05.999999999",
	"2006-01-02 15:04:05.999999999",
}

// formatLocalDateTime renders a wall-clock datetime as
// "2006-01-02T15:04:05[.fff]" with no timezone designator. Fractional
// seconds are emitted only when non-zero; trailing zeros are trimmed.
func formatLocalDateTime(ldt LocalDateTime) string {
	return time.Date(
		ldt.Year, ldt.Month, ldt.Day,
		ldt.Hour, ldt.Minute, ldt.Second, ldt.Nanosec,
		time.UTC,
	).Format("2006-01-02T15:04:05.999999999")
}

// parseLocalDateTime parses a wall-clock datetime in either the "T" form or
// the space-separated form. Inputs carrying a timezone designator are
// rejected with an error.
func parseLocalDateTime(s string) (LocalDateTime, error) {
	if hasOffsetSuffix(s) {
		return LocalDateTime{}, fmt.Errorf("LocalDateTime must not carry timezone info: %q", s)
	}
	for _, f := range localDateTimeFormats {
		parsed, err := time.Parse(f, s)
		if err == nil {
			return localDateTimeFromTime(parsed), nil
		}
	}
	return LocalDateTime{}, fmt.Errorf("cannot parse LocalDateTime from %q", s)
}

// localDateTimeFromTime extracts the wall-clock components of t and discards
// its location.
func localDateTimeFromTime(t time.Time) LocalDateTime {
	return LocalDateTime{
		Year:    t.Year(),
		Month:   t.Month(),
		Day:     t.Day(),
		Hour:    t.Hour(),
		Minute:  t.Minute(),
		Second:  t.Second(),
		Nanosec: t.Nanosecond(),
	}
}
