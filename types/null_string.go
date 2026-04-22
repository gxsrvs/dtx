package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullString is a nullable string. Valid reports whether Val holds a
// meaningful value (true) or a SQL/JSON NULL (false).
type NullString struct {
	Val   string
	Valid bool
}

// NewNullString constructs a valid NullString wrapping value. The empty
// string is preserved as a valid value; use NewNullStringEmpty or
// NSFromString if an empty input should collapse to NULL.
func NewNullString(value string) NullString {
	return NullString{Val: value, Valid: true}
}

// NewNullStringEmpty returns an invalid (NULL) NullString.
func NewNullStringEmpty() NullString {
	return NullString{Valid: false}
}

// NSFromString builds an sql.NullString from a raw string, treating the
// empty string as NULL. Handy when composing queries against the
// database/sql package directly.
//
//goland:noinspection GoMixedReceiverTypes
func NSFromString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

// IsEmpty reports whether the value is NULL (Valid == false).
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullString) IsEmpty() bool {
	return !thisVal.Valid
}

// IsZero reports whether the value is NULL (Valid == false). Mirroring
// time.Time.IsZero, this also enables encoding/json's `omitzero` tag
// (Go 1.24+) to elide invalid wrappers from marshalled output.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullString) IsZero() bool {
	return !thisVal.Valid
}

// GetNullString returns the underlying string of an sql.NullString, or ""
// when the value is NULL. Mirrors NSFromString on the read side.
//
//goland:noinspection GoMixedReceiverTypes
func GetNullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

// ToString returns the underlying string, or "" when NULL. Note that
// valid empty strings are indistinguishable from NULL in this output —
// use the Valid field directly to discriminate.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullString) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return thisVal.Val
}

// Value implements the database/sql/driver.Valuer interface. A NULL
// value is emitted as (nil, nil).
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullString) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

// Scan implements the database/sql.Scanner interface, delegating to
// sql.NullString so that the driver's NULL signalling is honoured.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullStringEmpty()
		return nil
	}
	*thisVal = NewNullString(s.String)
	return nil
}

// MarshalJSON renders the value as a JSON string, or null when empty.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullString) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJSON, nil
	}
	return json.Marshal(thisVal.Val)
}

// UnmarshalJSON parses a JSON string or the token null. Surrounding
// quotes and JSON escape sequences are decoded via encoding/json.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullString) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		thisVal.Val = ""
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	thisVal.Valid = true
	thisVal.Val = s
	return nil
}
