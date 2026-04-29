package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// LocalTime is a not-null wall-clock time-of-day without a timezone.
//
// Use LocalTime for values like "shop opens at 09:00" — a time-of-day with
// no offset. For points in time with a known offset, use OffsetTime instead.
//
// In(loc) and UTC() are not provided by design: LocalTime carries no
// timezone, so zone reinterpretation is meaningless. Use ToTime() to
// materialise a time.Time on the zero date in UTC.
type LocalTime struct {
	Hour    int
	Minute  int
	Second  int
	Nanosec int
}

// NewLocalTime constructs a LocalTime from its components.
func NewLocalTime(hour, minute, second, nanosec int) LocalTime {
	return LocalTime{Hour: hour, Minute: minute, Second: second, Nanosec: nanosec}
}

// LocalTimeFromTime extracts the time-of-day components of t and discards
// its date and location.
func LocalTimeFromTime(t time.Time) LocalTime {
	return localTimeFromTime(t)
}

// LocalTimeFromString parses a wall-clock time-of-day in ISO 8601 form
// ("HH:MM[:SS[.fff]]"). Inputs carrying a timezone designator are rejected.
func LocalTimeFromString(strValue string) (LocalTime, error) {
	if strValue == "" {
		return LocalTime{}, errors.New("cannot parse LocalTime from empty string")
	}
	return parseLocalTime(strValue)
}

// ToTime materialises the time-of-day on the zero date (year 0, January 1)
// in time.UTC.
func (thisVal LocalTime) ToTime() time.Time {
	return time.Date(0, 1, 1,
		thisVal.Hour, thisVal.Minute, thisVal.Second, thisVal.Nanosec,
		time.UTC,
	)
}

// ToString renders the value in the library's canonical format
// ("15:04:05[.fff]").
func (thisVal LocalTime) ToString() string {
	return formatLocalTime(thisVal)
}

// Value implements the database/sql/driver.Valuer interface. The value is
// emitted as a time.Time on the zero date in time.UTC.
func (thisVal LocalTime) Value() (driver.Value, error) {
	return thisVal.ToTime(), nil
}

// Scan implements the database/sql.Scanner interface. Accepts time.Time,
// string, or []byte. A NULL value is rejected — LocalTime cannot be empty.
func (thisVal *LocalTime) Scan(value interface{}) error {
	if value == nil {
		return errors.New("cannot scan NULL into LocalTime")
	}
	switch v := value.(type) {
	case time.Time:
		*thisVal = localTimeFromTime(v)
		return nil
	case string:
		parsed, err := parseLocalTime(v)
		if err != nil {
			return err
		}
		*thisVal = parsed
		return nil
	case []byte:
		parsed, err := parseLocalTime(string(v))
		if err != nil {
			return err
		}
		*thisVal = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into LocalTime", value)
	}
}

// MarshalJSON renders the value as a JSON string in the library's canonical
// format.
func (thisVal LocalTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(formatLocalTime(thisVal))
}

// UnmarshalJSON parses a JSON string in any of the formats accepted by
// LocalTimeFromString. JSON null is rejected — LocalTime cannot be empty.
func (thisVal *LocalTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" {
		return errors.New("cannot unmarshal null into LocalTime")
	}
	s := strings.Trim(sd, "\"")
	if s == "" {
		return errors.New("cannot unmarshal empty string into LocalTime")
	}
	parsed, err := parseLocalTime(s)
	if err != nil {
		return err
	}
	*thisVal = parsed
	return nil
}

// Before reports whether thisVal precedes other. Comparison operates
// on wall-clock components (hour/minute/second/nanosec) since
// LocalTime carries no timezone.
func (thisVal LocalTime) Before(other LocalTime) bool {
	return thisVal.ToTime().Before(other.ToTime())
}

// After reports whether thisVal is after other.
func (thisVal LocalTime) After(other LocalTime) bool {
	return thisVal.ToTime().After(other.ToTime())
}

// Equal reports whether thisVal and other have identical wall-clock
// components.
func (thisVal LocalTime) Equal(other LocalTime) bool {
	return thisVal == other
}

// Compare returns -1, 0, or +1 as thisVal is before, equal to, or
// after other.
func (thisVal LocalTime) Compare(other LocalTime) int {
	return thisVal.ToTime().Compare(other.ToTime())
}

// Add returns thisVal shifted by d. No modulo-24h enforcement: a
// duration that crosses midnight rotates the implicit date in the
// underlying time.Time, but the date is not part of LocalTime's
// serialised form. Consume with care.
func (thisVal LocalTime) Add(d time.Duration) LocalTime {
	return LocalTimeFromTime(thisVal.ToTime().Add(d))
}

// Sub returns thisVal − other as a time.Duration. The implicit dates
// participate in the calculation; for same-day pairs the result is
// sub-24h.
func (thisVal LocalTime) Sub(other LocalTime) time.Duration {
	return thisVal.ToTime().Sub(other.ToTime())
}

// Truncate returns thisVal rounded down to the nearest multiple of d
// since the zero time.
func (thisVal LocalTime) Truncate(d time.Duration) LocalTime {
	return LocalTimeFromTime(thisVal.ToTime().Truncate(d))
}
