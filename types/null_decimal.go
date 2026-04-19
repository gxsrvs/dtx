package types

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/shopspring/decimal"
)

// NullDecimal is a nullable arbitrary-precision decimal backed by
// github.com/shopspring/decimal. Valid reports whether Val holds a
// meaningful value (true) or a SQL/JSON NULL (false).
type NullDecimal struct {
	Val   decimal.Decimal
	Valid bool
}

// NewNullDecimal constructs a valid NullDecimal wrapping val.
func NewNullDecimal(val decimal.Decimal) NullDecimal {
	return NullDecimal{
		Val:   val,
		Valid: true,
	}
}

// NewNullDecimalEmpty returns an invalid (NULL) NullDecimal.
func NewNullDecimalEmpty() NullDecimal {
	return NullDecimal{
		Val:   decimal.Decimal{},
		Valid: false,
	}
}

// IsEmpty reports whether the value is NULL (Valid == false).
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDecimal) IsEmpty() bool {
	return !thisVal.Valid
}

// ToString renders the decimal via decimal.Decimal.String (no trailing
// zeros), or "" when NULL.
func (thisVal *NullDecimal) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return thisVal.Val.String()
}

// NullDecimalFromString parses a decimal from the string pointer. A nil
// pointer, an empty string, the tokens "null"/"nil" (case-insensitive),
// or a parse error all produce an invalid NullDecimal.
func NullDecimalFromString(strValue *string) NullDecimal {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullDecimalEmpty()
	}
	result, err := decimal.NewFromString(*strValue)
	if err != nil {
		return NewNullDecimalEmpty()
	}
	return NewNullDecimal(result)
}

// MulNullDecimals multiplies two NullDecimals. The result is NULL if
// either operand is NULL.
func MulNullDecimals(val1, val2 NullDecimal) NullDecimal {
	if !val1.Valid || !val2.Valid {
		return NewNullDecimalEmpty()
	}
	return NewNullDecimal(val1.Val.Mul(val2.Val))
}

// Scan implements the database/sql.Scanner interface. The decimal is
// received as text from the driver (via sql.NullString) and parsed
// through NullDecimalFromString so precision is preserved.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDecimal) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullDecimalEmpty()
		return nil
	}
	*thisVal = NullDecimalFromString(&s.String)
	return nil
}

// MarshalJSON renders the value as a JSON number (the format
// decimal.Decimal itself emits), or null when empty.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullDecimal) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(thisVal.Val)
}

// UnmarshalJSON parses a JSON number, JSON string, or the token null.
// Any other input is treated as a parse error.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDecimal) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		thisVal.Val = decimal.Decimal{}
		return nil
	}
	err := json.Unmarshal(data, &thisVal.Val)
	if err != nil {
		thisVal.Valid = false
		return err
	}
	thisVal.Valid = true
	return nil
}
