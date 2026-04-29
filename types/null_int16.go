package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// NullInt16 is a nullable int16. Valid reports whether Val holds a
// meaningful value (true) or a SQL/JSON NULL (false).
type NullInt16 struct {
	Val   int16
	Valid bool
}

// NewNullInt16 constructs a valid NullInt16 wrapping value.
func NewNullInt16(value int16) NullInt16 {
	return NullInt16{Val: value, Valid: true}
}

// NewNullInt16Empty returns an invalid (NULL) NullInt16.
func NewNullInt16Empty() NullInt16 {
	return NullInt16{Valid: false}
}

// NullInt16FromString parses an int16 from the string pointer. A nil
// pointer, an empty string, the tokens "null"/"nil" (case-insensitive),
// or a parse error all produce an invalid NullInt16.
func NullInt16FromString(strValue *string) NullInt16 {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullInt16Empty()
	}
	var result, err = strconv.ParseInt(*strValue, 10, 16)
	if err != nil {
		return NewNullInt16Empty()
	}
	return NewNullInt16(int16(result))
}

// NullInt16FromNullString parses an int16 from a NullString. An invalid
// or empty NullString, or a parse error, produce an invalid NullInt16.
func NullInt16FromNullString(str NullString) NullInt16 {
	if !str.Valid || str.Val == "" {
		return NewNullInt16Empty()
	}
	var result, err = strconv.ParseInt(str.Val, 10, 16)
	if err != nil {
		return NewNullInt16Empty()
	}
	return NewNullInt16(int16(result))
}

// IsEmpty reports whether the value is NULL (Valid == false).
func (thisVal *NullInt16) IsEmpty() bool {
	return !thisVal.Valid
}

// IsZero reports whether the value is NULL (Valid == false). Mirroring
// time.Time.IsZero, this also enables encoding/json's `omitzero` tag
// (Go 1.24+) to elide invalid wrappers from marshalled output.
func (thisVal NullInt16) IsZero() bool {
	return !thisVal.Valid
}

// ToString renders the value in decimal form, or "" when NULL.
func (thisVal NullInt16) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return fmt.Sprintf("%d", thisVal.Val)
}

// Value implements the database/sql/driver.Valuer interface. A NULL value
// is emitted as (nil, nil); the valid case is widened to int64 per the
// database/sql contract.
func (thisVal NullInt16) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return int64(thisVal.Val), nil
}

// Scan implements the database/sql.Scanner interface, delegating to
// sql.NullInt16 so that the driver's NULL signalling is honoured.
func (thisVal *NullInt16) Scan(value interface{}) error {
	var s sql.NullInt16
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullInt16Empty()
		return nil
	}
	*thisVal = NewNullInt16(s.Int16)
	return nil
}

// MarshalJSON renders the value as a JSON number, or null when empty.
// The valid path uses strconv.AppendInt directly to skip the reflection
// dispatch and intermediate buffer allocation that json.Marshal does
// for primitive numeric types.
func (thisVal NullInt16) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJSON, nil
	}
	// 8 bytes covers the maximum int16 width including a leading sign.
	return strconv.AppendInt(make([]byte, 0, 8), int64(thisVal.Val), 10), nil
}

// UnmarshalJSON parses a JSON number or the token null. Any other input
// is treated as a parse error.
func (thisVal *NullInt16) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		thisVal.Val = 0
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = int16(val)
	return nil
}
