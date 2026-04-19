package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

// NullLocalDateTime is the nullable counterpart of LocalDateTime. It
// shares the same serializer (formatLocalDateTime / parseLocalDateTime).
// Valid reports whether Val holds a meaningful value (true) or a
// SQL/JSON NULL (false).
type NullLocalDateTime struct {
	Val   LocalDateTime
	Valid bool
}

// NewNullLocalDateTime constructs a valid NullLocalDateTime wrapping value.
func NewNullLocalDateTime(value LocalDateTime) NullLocalDateTime {
	return NullLocalDateTime{Val: value, Valid: true}
}

// NewNullLocalDateTimeEmpty returns an invalid (NULL) NullLocalDateTime.
func NewNullLocalDateTimeEmpty() NullLocalDateTime {
	return NullLocalDateTime{Valid: false}
}

// NullLocalDateTimeFromString parses a local datetime from the string
// pointer. A nil pointer, an empty string, the tokens "null"/"nil"
// (case-insensitive), or a parse error all produce an invalid
// NullLocalDateTime.
func NullLocalDateTimeFromString(strValue *string) NullLocalDateTime {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullLocalDateTimeEmpty()
	}
	parsed, err := parseLocalDateTime(*strValue)
	if err != nil {
		return NewNullLocalDateTimeEmpty()
	}
	return NewNullLocalDateTime(parsed)
}

// IsEmpty reports whether the value is NULL (Valid == false).
func (thisVal *NullLocalDateTime) IsEmpty() bool {
	return !thisVal.Valid
}

// ToString renders the value in the library's canonical local datetime
// form, or "" when NULL.
func (thisVal *NullLocalDateTime) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatLocalDateTime(thisVal.Val)
}

// Value implements the database/sql/driver.Valuer interface. A NULL
// value is emitted as (nil, nil); the valid case is materialised as a
// time.Time in UTC (drivers interpret TIMESTAMP WITHOUT TIME ZONE as
// the wall-clock components).
func (thisVal NullLocalDateTime) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val.ToTime(time.UTC), nil
}

// Scan implements the database/sql.Scanner interface, delegating to
// sql.NullTime so that the driver's NULL signalling is honoured. The
// returned time's wall-clock components are preserved verbatim.
func (thisVal *NullLocalDateTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullLocalDateTimeEmpty()
		return nil
	}
	*thisVal = NewNullLocalDateTime(localDateTimeFromTime(s.Time))
	return nil
}

// MarshalJSON renders the value as a JSON string in canonical local
// datetime form, or null when empty.
func (thisVal NullLocalDateTime) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(formatLocalDateTime(thisVal.Val))
}

// UnmarshalJSON parses a JSON string containing a local datetime, or
// the token null. Any other input is treated as a parse error.
func (thisVal *NullLocalDateTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := parseLocalDateTime(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}
