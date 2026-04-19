package types

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

// NullUuid is a nullable UUID backed by github.com/google/uuid. Valid
// reports whether Val holds a meaningful value (true) or a SQL/JSON
// NULL (false).
type NullUuid struct {
	Val   uuid.UUID
	Valid bool
}

// NewNullUuid constructs a valid NullUuid wrapping val.
func NewNullUuid(val uuid.UUID) NullUuid {
	return NullUuid{
		Val:   val,
		Valid: true,
	}
}

// NewNullUuidEmpty returns an invalid (NULL) NullUuid.
func NewNullUuidEmpty() NullUuid {
	return NullUuid{
		Val:   uuid.UUID{},
		Valid: false,
	}
}

// IsEmpty reports whether the value is NULL (Valid == false).
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullUuid) IsEmpty() bool {
	return !thisVal.Valid
}

// ToString renders the UUID in canonical 8-4-4-4-12 hex form, or ""
// when NULL.
func (thisVal *NullUuid) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return thisVal.Val.String()
}

// NullUuidFromString parses a UUID from the string pointer. A nil
// pointer, an empty string, the tokens "null"/"nil" (case-insensitive),
// or a parse error all produce an invalid NullUuid.
func NullUuidFromString(strValue *string) NullUuid {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullUuidEmpty()
	}
	result, err := uuid.Parse(*strValue)
	if err != nil {
		return NewNullUuidEmpty()
	}
	return NewNullUuid(result)
}

// Scan implements the database/sql.Scanner interface. The UUID is
// received as text from the driver (via sql.NullString) and parsed
// through NullUuidFromString.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullUuid) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullUuidEmpty()
		return nil
	}
	*thisVal = NullUuidFromString(&s.String)
	return nil
}

// MarshalJSON renders the value as a JSON string in canonical UUID
// form, or null when empty.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullUuid) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(thisVal.Val)
}

// UnmarshalJSON parses a JSON string containing a UUID, or the token
// null. Any other input is treated as a parse error.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullUuid) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		thisVal.Val = uuid.UUID{}
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
