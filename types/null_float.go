package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// NullFloat is a nullable float64. Valid reports whether Val holds a
// meaningful value (true) or a SQL/JSON NULL (false).
type NullFloat struct {
	Val   float64
	Valid bool
}

// NewNullFloat constructs a valid NullFloat wrapping value.
func NewNullFloat(value float64) NullFloat {
	return NullFloat{Val: value, Valid: true}
}

// NewNullFloatEmpty returns an invalid (NULL) NullFloat.
func NewNullFloatEmpty() NullFloat {
	return NullFloat{Valid: false}
}

// NullFloatFromString parses a float64 from the string pointer. A nil
// pointer, an empty string, the tokens "null"/"nil" (case-insensitive),
// or a parse error all produce an invalid NullFloat.
func NullFloatFromString(strValue *string) NullFloat {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullFloatEmpty()
	}
	var result, err = strconv.ParseFloat(*strValue, 64)
	if err != nil {
		return NewNullFloatEmpty()
	}
	return NewNullFloat(result)
}

// IsEmpty reports whether the value is NULL (Valid == false).
func (thisVal *NullFloat) IsEmpty() bool {
	return !thisVal.Valid
}

// IsZero reports whether the value is NULL (Valid == false). Mirroring
// time.Time.IsZero, this also enables encoding/json's `omitzero` tag
// (Go 1.24+) to elide invalid wrappers from marshalled output.
func (thisVal NullFloat) IsZero() bool {
	return !thisVal.Valid
}

// ToString renders the value via fmt's %f verb, or "" when NULL.
func (thisVal NullFloat) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return fmt.Sprintf("%f", thisVal.Val)
}

// Value implements the database/sql/driver.Valuer interface. A NULL value
// is emitted as (nil, nil).
func (thisVal NullFloat) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

// Scan implements the database/sql.Scanner interface, delegating to
// sql.NullFloat64 so that the driver's NULL signalling is honoured.
func (thisVal *NullFloat) Scan(value interface{}) error {
	var s sql.NullFloat64
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullFloatEmpty()
		return nil
	}
	*thisVal = NewNullFloat(s.Float64)
	return nil
}

// MarshalJSON renders the value as a JSON number, or null when empty.
// The valid path uses strconv.AppendFloat directly to skip the
// reflection dispatch and intermediate buffer allocation that
// json.Marshal does for primitive numeric types. The 'g' verb mirrors
// encoding/json's own format for float64.
func (thisVal NullFloat) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJSON, nil
	}
	// 32 bytes is enough for any float64 in 'g' form including sign,
	// decimal point, and exponent.
	return strconv.AppendFloat(make([]byte, 0, 32), thisVal.Val, 'g', -1, 64), nil
}

// UnmarshalJSON parses a JSON number or the token null. Any other input
// is treated as a parse error.
func (thisVal *NullFloat) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		thisVal.Val = 0
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}
