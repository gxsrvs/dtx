package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// NullInt64 is a nullable int64. Valid reports whether Val holds a
// meaningful value (true) or a SQL/JSON NULL (false).
type NullInt64 struct {
	Val   int64
	Valid bool
}

// NewNullInt64 constructs a valid NullInt64 wrapping value.
func NewNullInt64(value int64) NullInt64 {
	return NullInt64{Val: value, Valid: true}
}

// NewNullInt64Empty returns an invalid (NULL) NullInt64.
func NewNullInt64Empty() NullInt64 {
	return NullInt64{Valid: false}
}

// NullInt64FromString parses an int64 from the string pointer. A nil
// pointer, an empty string, the tokens "null"/"nil" (case-insensitive),
// or a parse error all produce an invalid NullInt64.
func NullInt64FromString(strValue *string) NullInt64 {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullInt64Empty()
	}
	result, err := strconv.ParseInt(*strValue, 10, 64)
	if err != nil {
		return NewNullInt64Empty()
	}
	return NewNullInt64(result)
}

// NullInt64FromNullString parses an int64 from a NullString. An invalid
// or empty NullString, or a parse error, produce an invalid NullInt64.
func NullInt64FromNullString(str NullString) NullInt64 {
	if !str.Valid || str.Val == "" {
		return NewNullInt64Empty()
	}
	var result, err = strconv.ParseInt(str.Val, 10, 64)
	if err != nil {
		return NewNullInt64Empty()
	}
	return NewNullInt64(result)
}

// IsEmpty reports whether the value is NULL (Valid == false).
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt64) IsEmpty() bool {
	return !thisVal.Valid
}

// ToString renders the value in decimal form, or "" when NULL.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt64) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return fmt.Sprintf("%d", thisVal.Val)
}

// Value implements the database/sql/driver.Valuer interface. A NULL value
// is emitted as (nil, nil).
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullInt64) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

// Scan implements the database/sql.Scanner interface, delegating to
// sql.NullInt64 so that the driver's NULL signalling is honoured.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt64) Scan(value interface{}) error {
	var s sql.NullInt64
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullInt64Empty()
		return nil
	}
	*thisVal = NewNullInt64(s.Int64)
	return nil
}

// MarshalJSON renders the value as a JSON number, or null when empty.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullInt64) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(thisVal.Val)
}

// UnmarshalJSON parses a JSON number or the token null. Any other input
// is treated as a parse error.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt64) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		thisVal.Val = 0
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}
