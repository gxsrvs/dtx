package types

import "time"

// isDigitString reports whether every rune in s is an ASCII digit.
func isDigitString(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// hasOffsetSuffix reports whether s ends with an ISO 8601 timezone designator:
// "Z" / "z", "+HH:MM" / "-HH:MM", or "+HHMM" / "-HHMM". Used by the local
// (zone-less) parsers to reject inputs that carry a zone.
func hasOffsetSuffix(s string) bool {
	n := len(s)
	if n == 0 {
		return false
	}
	if c := s[n-1]; c == 'Z' || c == 'z' {
		return true
	}
	if n >= 6 {
		c := s[n-6]
		if (c == '+' || c == '-') &&
			isDigitString(s[n-5:n-3]) && s[n-3] == ':' && isDigitString(s[n-2:]) {
			return true
		}
	}
	if n >= 5 {
		c := s[n-5]
		if (c == '+' || c == '-') && isDigitString(s[n-4:]) {
			return true
		}
	}
	return false
}

func MaxDateTime(dt1, dt2 time.Time) time.Time {
	if dt1.After(dt2) {
		return dt1
	}
	return dt2
}

//goland:noinspection GoUnusedExportedFunction
func MinDateTime(dt1, dt2 time.Time) time.Time {
	if dt1.Before(dt2) {
		return dt1
	}
	return dt2
}

func AssembleDateTime(
	dateValue *time.Time,
	timeValue *time.Time,
	location *time.Location,
) *time.Time {
	y, m, d := dateValue.Date()
	hh, mm, ss, ns := timeValue.Hour(), timeValue.Minute(), timeValue.Second(), timeValue.Nanosecond()
	if location != nil {
		result := time.Date(
			y, m, d,
			hh, mm, ss, ns,
			location,
		)
		return &result
	}
	result := time.Date(
		y, m, d,
		hh, mm, ss, ns,
		timeValue.Location(),
	)
	return &result
}

func AssembleDateTimeTZ(
	dateValue *time.Time,
	timeValue *time.Time,
	timeZone string,
) (*time.Time, error) {
	loc, _, err := ParseTimezoneExtended(timeZone)
	if err != nil {
		return dateValue, err
	}
	return AssembleDateTime(dateValue, timeValue, loc), nil
}

func AssembleNullDateTimeTZ(
	dateValue *NullDate,
	defaultDate *time.Time,
	timeValue *NullOffsetTime,
	defaultTime *time.Time,
	timeZone string,
) (*time.Time, error) {
	var dateVal *time.Time
	if dateValue.Valid {
		dateVal = &dateValue.Val
	} else {
		dateVal = defaultDate
	}

	var timeVal *time.Time
	if timeValue.Valid {
		timeVal = &timeValue.Val
	} else {
		timeVal = defaultTime
	}

	return AssembleDateTimeTZ(dateVal, timeVal, timeZone)
}
