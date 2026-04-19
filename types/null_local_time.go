package types

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

// NullLocalTime is the nullable counterpart of LocalTime. It shares the
// same serializer (formatLocalTime / parseLocalTime). Valid reports
// whether Val holds a meaningful value (true) or a SQL/JSON NULL
// (false).
type NullLocalTime struct {
	Val   LocalTime
	Valid bool
}

// NewNullLocalTime constructs a valid NullLocalTime wrapping value.
func NewNullLocalTime(value LocalTime) NullLocalTime {
	return NullLocalTime{Val: value, Valid: true}
}

// NewNullLocalTimeEmpty returns an invalid (NULL) NullLocalTime.
func NewNullLocalTimeEmpty() NullLocalTime {
	return NullLocalTime{Valid: false}
}

// NullLocalTimeFromString parses a local time from the string pointer.
// A nil pointer, an empty string, the tokens "null"/"nil"
// (case-insensitive), or a parse error all produce an invalid
// NullLocalTime.
func NullLocalTimeFromString(strValue *string) NullLocalTime {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullLocalTimeEmpty()
	}
	parsed, err := parseLocalTime(*strValue)
	if err != nil {
		return NewNullLocalTimeEmpty()
	}
	return NewNullLocalTime(parsed)
}

// IsEmpty reports whether the value is NULL (Valid == false).
func (thisVal *NullLocalTime) IsEmpty() bool {
	return !thisVal.Valid
}

// ToString renders the value in the library's canonical local time
// form, or "" when NULL.
func (thisVal *NullLocalTime) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatLocalTime(thisVal.Val)
}

// Value implements the database/sql/driver.Valuer interface. A NULL
// value is emitted as (nil, nil); the valid case is materialised via
// LocalTime.ToTime (a same-day time.Time in UTC).
func (thisVal NullLocalTime) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val.ToTime(), nil
}

// Scan implements the database/sql.Scanner interface. NULL values
// produce an invalid NullLocalTime; non-NULL values are delegated to
// LocalTime.Scan.
func (thisVal *NullLocalTime) Scan(value interface{}) error {
	if value == nil {
		*thisVal = NewNullLocalTimeEmpty()
		return nil
	}
	var lt LocalTime
	if err := lt.Scan(value); err != nil {
		return err
	}
	*thisVal = NewNullLocalTime(lt)
	return nil
}

// MarshalJSON renders the value as a JSON string in canonical local
// time form, or null when empty.
func (thisVal NullLocalTime) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(formatLocalTime(thisVal.Val))
}

// UnmarshalJSON parses a JSON string containing a local time, or the
// token null. Any other input is treated as a parse error.
func (thisVal *NullLocalTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := parseLocalTime(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}
