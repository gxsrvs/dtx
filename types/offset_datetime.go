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

// Before reports whether the receiver is before dt.
func (thisVal OffsetDateTime) Before(dt time.Time) bool {
	return time.Time(thisVal).Before(dt)
}

// After reports whether the receiver is after dt.
func (thisVal OffsetDateTime) After(dt time.Time) bool {
	return time.Time(thisVal).After(dt)
}

// Unix returns the underlying instant as a Unix timestamp — the number
// of seconds elapsed since January 1, 1970 UTC.
func (thisVal OffsetDateTime) Unix() int64 {
	return time.Time(thisVal).Unix()
}
