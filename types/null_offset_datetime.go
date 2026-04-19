package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

// NullOffsetDateTime is the nullable counterpart of OffsetDateTime. It
// shares the same serializer (formatOffsetDateTime /
// parseOffsetDateTime). Valid reports whether Val holds a meaningful
// value (true) or a SQL/JSON NULL (false).
type NullOffsetDateTime struct {
	Val   time.Time
	Valid bool
}

// NewNullOffsetDateTime constructs a valid NullOffsetDateTime wrapping
// value.
func NewNullOffsetDateTime(value time.Time) NullOffsetDateTime {
	return NullOffsetDateTime{Val: value, Valid: true}
}

// NewNullOffsetDateTimeEmpty returns an invalid (NULL) NullOffsetDateTime.
func NewNullOffsetDateTimeEmpty() NullOffsetDateTime {
	return NullOffsetDateTime{Valid: false}
}

// NullOffsetDateTimeFromString parses an offset datetime from the
// string pointer. A nil pointer, an empty string, the tokens
// "null"/"nil" (case-insensitive), or a parse error all produce an
// invalid NullOffsetDateTime.
func NullOffsetDateTimeFromString(strValue *string) NullOffsetDateTime {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullOffsetDateTimeEmpty()
	}
	parsed, err := parseOffsetDateTime(*strValue)
	if err != nil {
		return NewNullOffsetDateTimeEmpty()
	}
	return NewNullOffsetDateTime(parsed)
}

// IsEmpty reports whether the value is NULL (Valid == false).
func (thisVal *NullOffsetDateTime) IsEmpty() bool {
	return !thisVal.Valid
}

// ToString renders the value in the library's canonical offset
// datetime form, or "" when NULL.
func (thisVal NullOffsetDateTime) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatOffsetDateTime(thisVal.Val)
}

// Value implements the database/sql/driver.Valuer interface. A NULL
// value is emitted as (nil, nil); the valid case is emitted as the
// underlying time.Time.
func (thisVal NullOffsetDateTime) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

// Scan implements the database/sql.Scanner interface, delegating to
// sql.NullTime so that the driver's NULL signalling is honoured.
func (thisVal *NullOffsetDateTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullOffsetDateTimeEmpty()
		return nil
	}
	*thisVal = NewNullOffsetDateTime(s.Time)
	return nil
}

// MarshalJSON renders the value as a JSON string in canonical offset
// datetime form, or null when empty.
func (thisVal NullOffsetDateTime) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(formatOffsetDateTime(thisVal.Val))
}

// UnmarshalJSON parses a JSON string containing an offset datetime, or
// the token null. Any other input is treated as a parse error.
func (thisVal *NullOffsetDateTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := parseOffsetDateTime(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}

// Before reports whether the receiver is valid and its instant is
// before dt. An invalid (NULL) value is never before anything.
func (thisVal NullOffsetDateTime) Before(dt time.Time) bool {
	return thisVal.Valid && thisVal.Val.Before(dt)
}

// After reports whether the receiver is valid and its instant is after
// dt. An invalid (NULL) value is never after anything.
func (thisVal NullOffsetDateTime) After(dt time.Time) bool {
	return thisVal.Valid && thisVal.Val.After(dt)
}
