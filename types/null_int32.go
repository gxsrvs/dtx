package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// NullInt32 is a nullable int32. Valid reports whether Val holds a
// meaningful value (true) or a SQL/JSON NULL (false).
type NullInt32 struct {
	Val   int32
	Valid bool
}

// NewNullInt32 constructs a valid NullInt32 wrapping value.
func NewNullInt32(value int32) NullInt32 {
	return NullInt32{Val: value, Valid: true}
}

// NewNullInt32Empty returns an invalid (NULL) NullInt32.
func NewNullInt32Empty() NullInt32 {
	return NullInt32{Valid: false}
}

// NullInt32FromString parses an int32 from the string pointer. A nil
// pointer, an empty string, the tokens "null"/"nil" (case-insensitive),
// or a parse error all produce an invalid NullInt32.
func NullInt32FromString(strValue *string) NullInt32 {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullInt32Empty()
	}
	var result, err = strconv.ParseInt(*strValue, 10, 32)
	if err != nil {
		return NewNullInt32Empty()
	}
	return NewNullInt32(int32(result))
}

// NullInt32FromNullString parses an int32 from a NullString. An invalid
// or empty NullString, or a parse error, produce an invalid NullInt32.
func NullInt32FromNullString(str NullString) NullInt32 {
	if !str.Valid || str.Val == "" {
		return NewNullInt32Empty()
	}
	var result, err = strconv.ParseInt(str.Val, 10, 32)
	if err != nil {
		return NewNullInt32Empty()
	}
	return NewNullInt32(int32(result))
}

// IsEmpty reports whether the value is NULL (Valid == false).
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt32) IsEmpty() bool {
	return !thisVal.Valid
}

// IsZero reports whether the value is NULL (Valid == false). Mirroring
// time.Time.IsZero, this also enables encoding/json's `omitzero` tag
// (Go 1.24+) to elide invalid wrappers from marshalled output.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullInt32) IsZero() bool {
	return !thisVal.Valid
}

// ToString renders the value in decimal form, or "" when NULL.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt32) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return fmt.Sprintf("%d", thisVal.Val)
}

// Value implements the database/sql/driver.Valuer interface. A NULL value
// is emitted as (nil, nil); the valid case is widened to int64 per the
// database/sql contract.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullInt32) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return int64(thisVal.Val), nil
}

// Scan implements the database/sql.Scanner interface, delegating to
// sql.NullInt32 so that the driver's NULL signalling is honoured.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt32) Scan(value interface{}) error {
	var s sql.NullInt32
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullInt32Empty()
		return nil
	}
	*thisVal = NewNullInt32(s.Int32)
	return nil
}

// MarshalJSON renders the value as a JSON number, or null when empty.
// The valid path uses strconv.AppendInt directly to skip the reflection
// dispatch and intermediate buffer allocation that json.Marshal does
// for primitive numeric types.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullInt32) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	// 12 bytes covers the maximum int32 width including a leading sign.
	return strconv.AppendInt(make([]byte, 0, 12), int64(thisVal.Val), 10), nil
}

// UnmarshalJSON parses a JSON number or the token null. Any other input
// is treated as a parse error.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt32) UnmarshalJSON(data []byte) error {
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
	thisVal.Val = int32(val)
	return nil
}
