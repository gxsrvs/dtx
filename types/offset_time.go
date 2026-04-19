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
