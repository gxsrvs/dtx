package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// OffsetTime is a not-null time-of-day with a timezone offset.
//
// Use OffsetTime when the time-of-day value carries a meaningful zone
// offset (e.g. exchange opens at "09:30:00-05:00"). For zone-less times,
// use LocalTime.
type OffsetTime time.Time

// NewOffsetTime wraps t as an OffsetTime, preserving its location.
func NewOffsetTime(t time.Time) OffsetTime {
	return OffsetTime(t)
}

// OffsetTimeFromString parses a time-of-day with an optional timezone
// designator. When no designator is present, time.Local is assumed.
func OffsetTimeFromString(strValue string) (OffsetTime, error) {
	if strValue == "" {
		return OffsetTime{}, errors.New("cannot parse OffsetTime from empty string")
	}
	parsed, err := parseOffsetTime(strValue)
	if err != nil {
		return OffsetTime{}, err
	}
	return OffsetTime(parsed), nil
}

// AsTime returns the underlying time.Time.
func (thisVal OffsetTime) AsTime() time.Time {
	return time.Time(thisVal)
}

// String renders the value in the library's canonical TZ-aware format.
func (thisVal OffsetTime) String() string {
	return formatOffsetTime(time.Time(thisVal))
}

// ToString is an alias for String, satisfying the ToStringAble interface.
func (thisVal OffsetTime) ToString() string {
	return formatOffsetTime(time.Time(thisVal))
}

// Value implements the database/sql/driver.Valuer interface.
func (thisVal OffsetTime) Value() (driver.Value, error) {
	return time.Time(thisVal), nil
}

// Scan implements the database/sql.Scanner interface. A NULL value is
// rejected — OffsetTime cannot be empty.
func (thisVal *OffsetTime) Scan(value interface{}) error {
	if value == nil {
		return errors.New("cannot scan NULL into OffsetTime")
	}
	switch v := value.(type) {
	case time.Time:
		*thisVal = OffsetTime(v)
		return nil
	case string:
		parsed, err := parseOffsetTime(v)
		if err != nil {
			return err
		}
		*thisVal = OffsetTime(parsed)
		return nil
	case []byte:
		parsed, err := parseOffsetTime(string(v))
		if err != nil {
			return err
		}
		*thisVal = OffsetTime(parsed)
		return nil
	default:
		return errors.New("unsupported type for OffsetTime.Scan")
	}
}

// MarshalJSON renders the value as a JSON string in the library's canonical
// TZ-aware format.
func (thisVal OffsetTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(formatOffsetTime(time.Time(thisVal)))
}

// UnmarshalJSON parses a JSON string in any of the formats accepted by
// OffsetTimeFromString. JSON null is rejected — OffsetTime cannot be empty.
func (thisVal *OffsetTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" {
		return errors.New("cannot unmarshal null into OffsetTime")
	}
	s := strings.Trim(sd, "\"")
	if s == "" {
		return errors.New("cannot unmarshal empty string into OffsetTime")
	}
	parsed, err := parseOffsetTime(s)
	if err != nil {
		return err
	}
	*thisVal = OffsetTime(parsed)
	return nil
}

// In returns the same instant rendered in loc. Mirrors time.Time.In:
// only the wall-clock representation changes. The implicit date hidden
// in the underlying time.Time may shift by ±1 day across the zone
// boundary, but the date is not part of the OffsetTime serialised form.
func (thisVal OffsetTime) In(loc *time.Location) OffsetTime {
	return OffsetTime(time.Time(thisVal).In(loc))
}

// UTC is shorthand for In(time.UTC).
func (thisVal OffsetTime) UTC() OffsetTime {
	return thisVal.In(time.UTC)
}

// Before reports whether thisVal precedes other.
func (thisVal OffsetTime) Before(other OffsetTime) bool {
	return time.Time(thisVal).Before(time.Time(other))
}

// After reports whether thisVal is after other.
func (thisVal OffsetTime) After(other OffsetTime) bool {
	return time.Time(thisVal).After(time.Time(other))
}

// Equal reports whether thisVal and other denote the same instant.
// Two values that differ only in zone (e.g. "10:00:00+04:00" and
// "06:00:00Z") are Equal.
func (thisVal OffsetTime) Equal(other OffsetTime) bool {
	return time.Time(thisVal).Equal(time.Time(other))
}

// Compare returns -1, 0, or +1 as thisVal is before, equal to, or
// after other.
func (thisVal OffsetTime) Compare(other OffsetTime) int {
	return time.Time(thisVal).Compare(time.Time(other))
}

// Add returns thisVal shifted by d. No modulo-24h enforcement: a
// duration that crosses midnight rotates the implicit date in the
// underlying time.Time, but the date is not part of OffsetTime's
// serialised form. Consume with care.
func (thisVal OffsetTime) Add(d time.Duration) OffsetTime {
	return OffsetTime(time.Time(thisVal).Add(d))
}

// Sub returns thisVal − other as a time.Duration. The implicit dates
// participate in the calculation (so a same-day pair yields a
// sub-24h difference), but Durations larger than 24h are possible if
// the operands were constructed at different days.
func (thisVal OffsetTime) Sub(other OffsetTime) time.Duration {
	return time.Time(thisVal).Sub(time.Time(other))
}

// Truncate returns thisVal rounded down to the nearest multiple of d
// since the zero time.
func (thisVal OffsetTime) Truncate(d time.Duration) OffsetTime {
	return OffsetTime(time.Time(thisVal).Truncate(d))
}

// Hour returns the hour of the day (0-23).
func (thisVal OffsetTime) Hour() int { return time.Time(thisVal).Hour() }

// Minute returns the minute of the hour (0-59).
func (thisVal OffsetTime) Minute() int { return time.Time(thisVal).Minute() }

// Second returns the second of the minute (0-59).
func (thisVal OffsetTime) Second() int { return time.Time(thisVal).Second() }

// Nanosecond returns the nanosecond offset within the second
// (0-999999999).
func (thisVal OffsetTime) Nanosecond() int { return time.Time(thisVal).Nanosecond() }
