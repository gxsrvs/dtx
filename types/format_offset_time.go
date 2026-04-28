package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// offsetTimeFormat is the canonical TZ-aware time-of-day layout. The
// "Z07:00" suffix makes time.Format emit "Z" when the offset is zero and
// "+HH:MM" / "-HH:MM" otherwise; ".999999999" trims trailing zeros and
// makes fractional seconds optional on parse.
const offsetTimeFormat = "15:04:05.999999999Z07:00"

// timeFormats lists the time-of-day layouts accepted by ParseTimeFromString,
// in the order they are tried. The timezone designator (if any) is stripped
// by ParseTimezoneExtended before reaching this list.
var timeFormats = []string{
	"15:04",
	"15:04:05",
	"15:04:05.000",
	"15:04:05.000000",
	"15:04:05.000000000",
}

// formatOffsetTime renders a time-of-day in the library's canonical TZ-aware
// format ("15:04:05[.fff…]" with "Z" or "±HH:MM" suffix).
func formatOffsetTime(t time.Time) string {
	return t.Format(offsetTimeFormat)
}

// parseOffsetTime parses a time-of-day with optional timezone designator.
// It is the value-returning counterpart of the public ParseTimeFromString.
func parseOffsetTime(s string) (time.Time, error) {
	parsed, err := ParseTimeFromString(s)
	if err != nil {
		return time.Time{}, err
	}
	return *parsed, nil
}

// ParseTimeFromString parses a time-of-day with an optional timezone
// designator ("HH:MM[:SS[.fff…]][Z|±HH:MM|±HHMM]").
func ParseTimeFromString(strValue string) (*time.Time, error) {
	loc, strValueRest, err := ParseTimezoneExtended(strValue)
	if err != nil {
		return nil, err
	}

	strValueRest = strings.TrimSpace(strValueRest)

	var parsed time.Time
	for _, f := range timeFormats {
		parsed, err = time.ParseInLocation(f, strValueRest, loc)
		if err == nil {
			return &parsed, nil
		}
	}

	return nil, errors.New("cannot parse time from " + strValue)
}

// ParseTimezoneExtended attempts to extract a trailing timezone designator
// from strValue. Supported forms are "Z" (RFC 3339 UTC literal),
// "+HH:MM" / "-HH:MM" (six chars) and "+HHMM" / "-HHMM" (five chars). When
// a designator is found it is removed from the returned string and
// represented as a *time.Location; otherwise time.Local is returned.
func ParseTimezoneExtended(strValue string) (*time.Location, string, error) {
	loc := time.Local
	s := strValue

	if len(s) >= 1 && s[len(s)-1] == 'Z' {
		return time.UTC, s[:len(s)-1], nil
	}

	if len(s) >= 6 {
		tzCandidate := s[len(s)-6:]
		if (tzCandidate[0] == '+' || tzCandidate[0] == '-' || tzCandidate[0] == ' ') &&
			isDigitString(tzCandidate[1:3]) && tzCandidate[3] == ':' &&
			isDigitString(tzCandidate[4:]) {

			offsetLoc, err := parseTimezoneWithColon(tzCandidate)
			if err == nil {
				loc = offsetLoc
				s = s[:len(s)-6]
				return loc, s, nil
			}
		}
	}

	if len(s) >= 5 {
		tzCandidate := s[len(s)-5:]
		if (tzCandidate[0] == '+' || tzCandidate[0] == '-') &&
			isDigitString(tzCandidate[1:]) {

			offsetLoc, err := parseTimezone(tzCandidate)
			if err == nil {
				loc = offsetLoc
				s = s[:len(s)-5]
				return loc, s, nil
			}
		}
	}

	return loc, s, nil
}

// parseTimezoneWithColon parses a "+HH:MM" / "-HH:MM" timezone designator.
func parseTimezoneWithColon(tz string) (*time.Location, error) {
	sign := tz[0]
	hhStr := tz[1:3]
	mmStr := tz[4:6]

	hh, err := strconv.Atoi(hhStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse hours from timezone %s", tz)
	}
	mm, err := strconv.Atoi(mmStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse minutes from timezone %s", tz)
	}

	offset := hh*3600 + mm*60
	if sign == '-' {
		offset = -offset
	}
	return time.FixedZone(fmt.Sprintf("UTC%s", tz), offset), nil
}

// parseTimezone parses a "+HHMM" / "-HHMM" timezone designator.
func parseTimezone(tz string) (*time.Location, error) {
	if len(tz) != 5 {
		return nil, errors.New("unknown time zone format")
	}
	sign := tz[0]
	hhStr := tz[1:3]
	mmStr := tz[3:5]

	hh, err := strconv.Atoi(hhStr)
	if err != nil {
		return nil, fmt.Errorf("timezone hours parsing error: %w", err)
	}
	mm, err := strconv.Atoi(mmStr)
	if err != nil {
		return nil, fmt.Errorf("timezone minutes parsing error: %w", err)
	}

	offset := hh*3600 + mm*60
	if sign == '-' {
		offset = -offset
	}
	loc := time.FixedZone(fmt.Sprintf("UTC%s", tz), offset)
	return loc, nil
}
