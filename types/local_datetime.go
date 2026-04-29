package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// LocalDateTime is a not-null wall-clock datetime without a timezone.
//
// Use LocalDateTime for business events anchored to a calendar moment in a
// human's wall-clock time (e.g. "contract signed on 2026-04-18 at 13:00")
// where the offset is intentionally absent. For points in time with a known
// offset, use OffsetDateTime instead.
//
// In(loc) and UTC() are not provided by design: LocalDateTime carries no
// timezone, so zone reinterpretation is meaningless. Use ToTime(loc) to
// materialise a time.Time in a chosen zone.
type LocalDateTime struct {
	Year    int
	Month   time.Month
	Day     int
	Hour    int
	Minute  int
	Second  int
	Nanosec int
}

// NewLocalDateTime constructs a LocalDateTime from its components.
func NewLocalDateTime(year int, month time.Month, day, hour, minute, second, nanosec int) LocalDateTime {
	return LocalDateTime{
		Year:    year,
		Month:   month,
		Day:     day,
		Hour:    hour,
		Minute:  minute,
		Second:  second,
		Nanosec: nanosec,
	}
}

// LocalDateTimeFromTime extracts the wall-clock components of t and discards
// its location.
func LocalDateTimeFromTime(t time.Time) LocalDateTime {
	return localDateTimeFromTime(t)
}

// LocalDateTimeFromString parses a wall-clock datetime in either ISO 8601
// "T" form ("2006-01-02T15:04:05[.fff]") or the equivalent space-separated
// form. Inputs carrying a timezone designator are rejected.
func LocalDateTimeFromString(strValue string) (LocalDateTime, error) {
	if strValue == "" {
		return LocalDateTime{}, errors.New("cannot parse LocalDateTime from empty string")
	}
	return parseLocalDateTime(strValue)
}

// ToTime materialises the wall-clock components in the given location. Pass
// time.UTC, time.Local, or a fixed zone according to how the value should be
// interpreted at the boundary with timezone-aware code. A nil loc is treated
// as time.UTC.
func (thisVal LocalDateTime) ToTime(loc *time.Location) time.Time {
	if loc == nil {
		loc = time.UTC
	}
	return time.Date(
		thisVal.Year, thisVal.Month, thisVal.Day,
		thisVal.Hour, thisVal.Minute, thisVal.Second, thisVal.Nanosec,
		loc,
	)
}

// ToString renders the value in the library's canonical format
// ("2006-01-02T15:04:05[.fff]").
func (thisVal LocalDateTime) ToString() string {
	return formatLocalDateTime(thisVal)
}

// Value implements the database/sql/driver.Valuer interface. The value is
// emitted in time.UTC; the SQL column type should be TIMESTAMP WITHOUT TIME
// ZONE so the driver writes the wall-clock fields verbatim.
func (thisVal LocalDateTime) Value() (driver.Value, error) {
	return thisVal.ToTime(time.UTC), nil
}

// Scan implements the database/sql.Scanner interface. A NULL value is
// rejected — LocalDateTime cannot be empty.
func (thisVal *LocalDateTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		return errors.New("cannot scan NULL into LocalDateTime")
	}
	*thisVal = localDateTimeFromTime(s.Time)
	return nil
}

// MarshalJSON renders the value as a JSON string in the library's canonical
// format.
func (thisVal LocalDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(formatLocalDateTime(thisVal))
}

// UnmarshalJSON parses a JSON string in any of the formats accepted by
// LocalDateTimeFromString. JSON null is rejected — LocalDateTime cannot
// be empty.
func (thisVal *LocalDateTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" {
		return errors.New("cannot unmarshal null into LocalDateTime")
	}
	s := strings.Trim(sd, "\"")
	if s == "" {
		return errors.New("cannot unmarshal empty string into LocalDateTime")
	}
	parsed, err := parseLocalDateTime(s)
	if err != nil {
		return err
	}
	*thisVal = parsed
	return nil
}

// Before reports whether thisVal precedes other. Comparison operates
// on wall-clock components (year/month/.../nanosec) since LocalDateTime
// carries no timezone.
func (thisVal LocalDateTime) Before(other LocalDateTime) bool {
	return thisVal.ToTime(time.UTC).Before(other.ToTime(time.UTC))
}

// After reports whether thisVal is after other.
func (thisVal LocalDateTime) After(other LocalDateTime) bool {
	return thisVal.ToTime(time.UTC).After(other.ToTime(time.UTC))
}

// Equal reports whether thisVal and other have identical wall-clock
// components.
func (thisVal LocalDateTime) Equal(other LocalDateTime) bool {
	return thisVal == other
}

// Compare returns -1, 0, or +1 as thisVal is before, equal to, or
// after other.
func (thisVal LocalDateTime) Compare(other LocalDateTime) int {
	return thisVal.ToTime(time.UTC).Compare(other.ToTime(time.UTC))
}

// Add returns thisVal shifted by d on the wall-clock components.
func (thisVal LocalDateTime) Add(d time.Duration) LocalDateTime {
	return LocalDateTimeFromTime(thisVal.ToTime(time.UTC).Add(d))
}

// AddDate returns thisVal with the given years/months/days added.
func (thisVal LocalDateTime) AddDate(years, months, days int) LocalDateTime {
	return LocalDateTimeFromTime(thisVal.ToTime(time.UTC).AddDate(years, months, days))
}

// Sub returns thisVal − other as a time.Duration on the wall-clock
// components.
func (thisVal LocalDateTime) Sub(other LocalDateTime) time.Duration {
	return thisVal.ToTime(time.UTC).Sub(other.ToTime(time.UTC))
}

// Truncate returns thisVal rounded down to the nearest multiple of d
// since the zero time.
func (thisVal LocalDateTime) Truncate(d time.Duration) LocalDateTime {
	return LocalDateTimeFromTime(thisVal.ToTime(time.UTC).Truncate(d))
}

// YearDay returns the day-of-year (1-365 / 1-366 in leap years).
// Year/Month/Day/Hour/Minute/Second/Nanosec are exposed as struct
// fields directly.
func (thisVal LocalDateTime) YearDay() int {
	return thisVal.ToTime(time.UTC).YearDay()
}
