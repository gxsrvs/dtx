package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

// NullOffsetTime is the nullable counterpart of OffsetTime. It shares
// the same serializer (formatOffsetTime / parseOffsetTime). Valid
// reports whether Val holds a meaningful value (true) or a SQL/JSON
// NULL (false).
type NullOffsetTime struct {
	Val   time.Time
	Valid bool
}

// NewNullOffsetTime constructs a valid NullOffsetTime wrapping value.
func NewNullOffsetTime(value time.Time) NullOffsetTime {
	return NullOffsetTime{Val: value, Valid: true}
}

// NewNullOffsetTimeEmpty returns an invalid (NULL) NullOffsetTime.
func NewNullOffsetTimeEmpty() NullOffsetTime {
	return NullOffsetTime{Valid: false}
}

// NullOffsetTimeFromString parses an offset time from the string
// pointer. A nil pointer, an empty string, the tokens "null"/"nil"
// (case-insensitive), or a parse error all produce an invalid
// NullOffsetTime.
func NullOffsetTimeFromString(strValue *string) NullOffsetTime {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullOffsetTimeEmpty()
	}
	parsed, err := parseOffsetTime(*strValue)
	if err != nil {
		return NewNullOffsetTimeEmpty()
	}
	return NewNullOffsetTime(parsed)
}

// IsEmpty reports whether the value is NULL (Valid == false).
func (thisVal *NullOffsetTime) IsEmpty() bool {
	return !thisVal.Valid
}

// ToString renders the value in the library's canonical offset time
// form, or "" when NULL.
func (thisVal NullOffsetTime) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatOffsetTime(thisVal.Val)
}

// Value implements the database/sql/driver.Valuer interface. A NULL
// value is emitted as (nil, nil); the valid case is emitted as the
// underlying time.Time.
func (thisVal NullOffsetTime) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

// Scan implements the database/sql.Scanner interface, delegating to
// sql.NullTime so that the driver's NULL signalling is honoured.
func (thisVal *NullOffsetTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullOffsetTimeEmpty()
		return nil
	}
	*thisVal = NewNullOffsetTime(s.Time)
	return nil
}

// MarshalJSON renders the value as a JSON string in canonical offset
// time form, or null when empty.
func (thisVal NullOffsetTime) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(formatOffsetTime(thisVal.Val))
}

// UnmarshalJSON parses a JSON string containing an offset time, or the
// token null. Any other input is treated as a parse error.
func (thisVal *NullOffsetTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := parseOffsetTime(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}
