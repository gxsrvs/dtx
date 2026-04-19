package types

import (
	"errors"
	"time"
)

// offsetDateTimeFormat is the canonical TZ-aware datetime layout
// (RFC 3339 with optional sub-second precision). The "Z07:00" suffix makes
// time.Format emit "Z" when the offset is zero and "+HH:MM" / "-HH:MM"
// otherwise; ".999999999" trims trailing zeros and keeps fractional seconds
// optional on parse.
const offsetDateTimeFormat = "2006-01-02T15:04:05.999999999Z07:00"

// formatOffsetDateTime renders a datetime in the library's canonical
// TZ-aware format.
func formatOffsetDateTime(t time.Time) string {
	return t.Format(offsetDateTimeFormat)
}

// parseOffsetDateTime parses a datetime with optional timezone designator.
// It is the value-returning counterpart of the public
// ParseDateTimeFromString.
func parseOffsetDateTime(s string) (time.Time, error) {
	parsed, err := ParseDateTimeFromString(s)
	if err != nil {
		return time.Time{}, err
	}
	return *parsed, nil
}

// ParseDateTimeFromString parses a datetime string. The date and time parts
// must be separated by either 'T' or ' '. The time part may carry an
// optional timezone designator ("Z", "+HH:MM", "+HHMM").
func ParseDateTimeFromString(strValue string) (*time.Time, error) {
	if len(strValue) < 11 {
		return nil, errors.New("cannot parse time from " + strValue)
	}
	sep := strValue[10]
	if sep != ' ' && sep != 'T' {
		return nil, errors.New("cannot parse time from " + strValue)
	}

	datePart := strValue[:10]
	parsedDate, err := ParseDateFromString(datePart)
	if err != nil {
		return nil, err
	}

	timePart := strValue[11:]
	parsedTime, err := ParseTimeFromString(timePart)
	if err != nil {
		return nil, err
	}

	result := AssembleDateTime(parsedDate, parsedTime, nil)
	return result, nil
}

// DateTimeToString renders a datetime in the library's canonical TZ-aware
// format. Kept as a thin wrapper for callers that operate on bare
// time.Time values.
func DateTimeToString(t time.Time) string {
	return formatOffsetDateTime(t)
}
