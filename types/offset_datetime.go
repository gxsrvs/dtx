package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// OffsetDateTime is a not-null datetime with a timezone offset (RFC 3339).
//
// Use OffsetDateTime when the value is anchored to a specific point in time
// (e.g. an event timestamp). For wall-clock datetimes without an offset,
// use LocalDateTime.
type OffsetDateTime time.Time

// NewOffsetDateTime wraps t as an OffsetDateTime, preserving its location.
func NewOffsetDateTime(t time.Time) OffsetDateTime {
	return OffsetDateTime(t)
}

// OffsetDateTimeFromString parses a datetime in any of the formats accepted
// by ParseDateTimeFromString.
func OffsetDateTimeFromString(strValue string) (OffsetDateTime, error) {
	if strValue == "" {
		return OffsetDateTime{}, errors.New("cannot parse OffsetDateTime from empty string")
	}
	parsed, err := parseOffsetDateTime(strValue)
	if err != nil {
		return OffsetDateTime{}, err
	}
	return OffsetDateTime(parsed), nil
}

// AsTime returns the underlying time.Time.
func (thisVal OffsetDateTime) AsTime() time.Time {
	return time.Time(thisVal)
}

// ToString renders the value in the library's canonical TZ-aware format.
func (thisVal OffsetDateTime) ToString() string {
	return formatOffsetDateTime(time.Time(thisVal))
}

// Value implements the database/sql/driver.Valuer interface.
func (thisVal OffsetDateTime) Value() (driver.Value, error) {
	return time.Time(thisVal), nil
}

// Scan implements the database/sql.Scanner interface. A NULL value is
// rejected — OffsetDateTime cannot be empty.
func (thisVal *OffsetDateTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		return errors.New("cannot scan NULL into OffsetDateTime")
	}
	*thisVal = OffsetDateTime(s.Time)
	return nil
}

// MarshalJSON renders the value as a JSON string in the library's canonical
// TZ-aware format.
func (thisVal OffsetDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(formatOffsetDateTime(time.Time(thisVal)))
}

// UnmarshalJSON parses a JSON string in any of the formats accepted by
// OffsetDateTimeFromString. JSON null is rejected — OffsetDateTime cannot
// be empty.
func (thisVal *OffsetDateTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" {
		return errors.New("cannot unmarshal null into OffsetDateTime")
	}
	s := strings.Trim(sd, "\"")
	if s == "" {
		return errors.New("cannot unmarshal empty string into OffsetDateTime")
	}
	parsed, err := parseOffsetDateTime(s)
	if err != nil {
		return err
	}
	*thisVal = OffsetDateTime(parsed)
	return nil
}

// Before reports whether thisVal precedes other.
func (thisVal OffsetDateTime) Before(other OffsetDateTime) bool {
	return time.Time(thisVal).Before(time.Time(other))
}

// After reports whether thisVal is after other.
func (thisVal OffsetDateTime) After(other OffsetDateTime) bool {
	return time.Time(thisVal).After(time.Time(other))
}

// Equal reports whether thisVal and other denote the same instant.
// Two values that differ only in zone (e.g. "10:00:00+04:00" and
// "06:00:00Z") are Equal.
func (thisVal OffsetDateTime) Equal(other OffsetDateTime) bool {
	return time.Time(thisVal).Equal(time.Time(other))
}

// Compare returns -1, 0, or +1 as thisVal is before, equal to, or
// after other.
func (thisVal OffsetDateTime) Compare(other OffsetDateTime) int {
	return time.Time(thisVal).Compare(time.Time(other))
}

// Add returns thisVal shifted by d.
func (thisVal OffsetDateTime) Add(d time.Duration) OffsetDateTime {
	return OffsetDateTime(time.Time(thisVal).Add(d))
}

// AddDate returns thisVal with the given years/months/days added.
func (thisVal OffsetDateTime) AddDate(years, months, days int) OffsetDateTime {
	return OffsetDateTime(time.Time(thisVal).AddDate(years, months, days))
}

// Sub returns thisVal − other as a time.Duration.
func (thisVal OffsetDateTime) Sub(other OffsetDateTime) time.Duration {
	return time.Time(thisVal).Sub(time.Time(other))
}

// Truncate returns thisVal rounded down to the nearest multiple of d
// since the zero time. Mirrors time.Time.Truncate.
func (thisVal OffsetDateTime) Truncate(d time.Duration) OffsetDateTime {
	return OffsetDateTime(time.Time(thisVal).Truncate(d))
}

// Unix returns the underlying instant as a Unix timestamp — the number
// of seconds elapsed since January 1, 1970 UTC.
func (thisVal OffsetDateTime) Unix() int64 {
	return time.Time(thisVal).Unix()
}

// UnixMilli returns the underlying instant as milliseconds since
// 1970-01-01 UTC.
func (thisVal OffsetDateTime) UnixMilli() int64 {
	return time.Time(thisVal).UnixMilli()
}

// UnixMicro returns the underlying instant as microseconds since
// 1970-01-01 UTC.
func (thisVal OffsetDateTime) UnixMicro() int64 {
	return time.Time(thisVal).UnixMicro()
}

// UnixNano returns the underlying instant as nanoseconds since
// 1970-01-01 UTC.
func (thisVal OffsetDateTime) UnixNano() int64 {
	return time.Time(thisVal).UnixNano()
}

// Year returns the calendar year (in the value's zone).
func (thisVal OffsetDateTime) Year() int { return time.Time(thisVal).Year() }

// Month returns the calendar month.
func (thisVal OffsetDateTime) Month() time.Month { return time.Time(thisVal).Month() }

// Day returns the day of the month (1-31).
func (thisVal OffsetDateTime) Day() int { return time.Time(thisVal).Day() }

// YearDay returns the day-of-year (1-365 / 1-366 in leap years).
func (thisVal OffsetDateTime) YearDay() int { return time.Time(thisVal).YearDay() }

// Hour returns the hour of the day (0-23).
func (thisVal OffsetDateTime) Hour() int { return time.Time(thisVal).Hour() }

// Minute returns the minute of the hour (0-59).
func (thisVal OffsetDateTime) Minute() int { return time.Time(thisVal).Minute() }

// Second returns the second of the minute (0-59).
func (thisVal OffsetDateTime) Second() int { return time.Time(thisVal).Second() }

// Nanosecond returns the nanosecond offset within the second
// (0-999999999).
func (thisVal OffsetDateTime) Nanosecond() int { return time.Time(thisVal).Nanosecond() }

// In returns the same instant rendered in loc. Mirrors time.Time.In:
// the absolute moment is preserved, only the displayed offset / wall
// clock changes (e.g. "10:00:00+04:00" becomes "06:00:00Z" when
// converted to UTC).
func (thisVal OffsetDateTime) In(loc *time.Location) OffsetDateTime {
	return OffsetDateTime(time.Time(thisVal).In(loc))
}

// UTC is shorthand for In(time.UTC).
func (thisVal OffsetDateTime) UTC() OffsetDateTime {
	return thisVal.In(time.UTC)
}
