package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// NullBool is a nullable boolean. Valid reports whether Val holds a
// meaningful value (true) or a SQL/JSON NULL (false).
type NullBool struct {
	Val   bool
	Valid bool
}

// NewNullBool constructs a valid NullBool wrapping value.
func NewNullBool(value bool) NullBool {
	return NullBool{Val: value, Valid: true}
}

// NewNullBoolEmpty returns an invalid (NULL) NullBool.
func NewNullBoolEmpty() NullBool {
	return NullBool{Valid: false}
}

// NullBoolFromString parses a bool from the string pointer. A nil pointer,
// an empty string, the tokens "null"/"nil" (case-insensitive), or a parse
// error all produce an invalid NullBool.
func NullBoolFromString(strValue *string) NullBool {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullBoolEmpty()
	}
	var result, err = strconv.ParseBool(*strValue)
	if err != nil {
		return NewNullBoolEmpty()
	}
	return NewNullBool(result)
}

// IsEmpty reports whether the value is NULL (Valid == false).
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullBool) IsEmpty() bool {
	return !thisVal.Valid
}

// ToString renders the value as "true"/"false", or "" when NULL.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullBool) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return fmt.Sprintf("%t", thisVal.Val)
}

// Value implements the database/sql/driver.Valuer interface. A NULL value
// is emitted as (nil, nil).
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullBool) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

// Scan implements the database/sql.Scanner interface, delegating to
// sql.NullBool so that the driver's NULL signalling is honoured.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullBool) Scan(value interface{}) error {
	var s sql.NullBool
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullBoolEmpty()
		return nil
	}
	*thisVal = NewNullBool(s.Bool)
	return nil
}

// MarshalJSON renders the value as a JSON boolean, or null when empty.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullBool) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(thisVal.Val)
}

// UnmarshalJSON parses a JSON boolean or the token null. Any other input
// is treated as a parse error.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullBool) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		thisVal.Val = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}
