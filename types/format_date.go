package types

import (
	"errors"
	"time"
)

// RuOnlyDateMask is the dd.MM.yyyy locale date format accepted by parseDate
// in addition to ISO 8601 (time.DateOnly).
const RuOnlyDateMask = "02.01.2006"

var dateFormats = []string{
	time.DateOnly,
	RuOnlyDateMask,
}

// formatDate renders a calendar date as "2006-01-02".
// It is the single source of truth for date marshalling shared by Date and
// NullDate.
func formatDate(t time.Time) string {
	return t.Format(time.DateOnly)
}

// parseDate parses a calendar date string using the formats listed in
// dateFormats. It is the single source of truth for date parsing shared by
// Date and NullDate.
func parseDate(s string) (time.Time, error) {
	for _, f := range dateFormats {
		parsed, err := time.ParseInLocation(f, s, time.Local)
		if err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, errors.New("cannot parse date from " + s)
}
