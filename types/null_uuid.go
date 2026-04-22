package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

// NullUUID is a nullable UUID backed by github.com/google/uuid. Valid
// reports whether Val holds a meaningful value (true) or a SQL/JSON
// NULL (false).
type NullUUID struct {
	Val   uuid.UUID
	Valid bool
}

// NewNullUUID constructs a valid NullUUID wrapping val.
func NewNullUUID(val uuid.UUID) NullUUID {
	return NullUUID{
		Val:   val,
		Valid: true,
	}
}

// NewNullUUIDEmpty returns an invalid (NULL) NullUUID.
func NewNullUUIDEmpty() NullUUID {
	return NullUUID{
		Val:   uuid.UUID{},
		Valid: false,
	}
}

// IsEmpty reports whether the value is NULL (Valid == false).
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullUUID) IsEmpty() bool {
	return !thisVal.Valid
}

// IsZero reports whether the value is NULL (Valid == false). Mirroring
// time.Time.IsZero, this also enables encoding/json's `omitzero` tag
// (Go 1.24+) to elide invalid wrappers from marshalled output.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullUUID) IsZero() bool {
	return !thisVal.Valid
}

// ToString renders the UUID in canonical 8-4-4-4-12 hex form, or ""
// when NULL.
func (thisVal *NullUUID) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return thisVal.Val.String()
}

// NullUUIDFromString parses a UUID from the string pointer. A nil
// pointer, an empty string, the tokens "null"/"nil" (case-insensitive),
// or a parse error all produce an invalid NullUUID.
func NullUUIDFromString(strValue *string) NullUUID {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullUUIDEmpty()
	}
	result, err := uuid.Parse(*strValue)
	if err != nil {
		return NewNullUUIDEmpty()
	}
	return NewNullUUID(result)
}

// Value implements the database/sql/driver.Valuer interface. A NULL
// value is emitted as (nil, nil); the valid case delegates to
// uuid.UUID.Value, which renders the UUID as its canonical
// 8-4-4-4-12 hex string.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullUUID) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val.Value()
}

// Scan implements the database/sql.Scanner interface. The UUID is
// received as text from the driver (via sql.NullString) and parsed
// through NullUUIDFromString.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullUUID) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullUUIDEmpty()
		return nil
	}
	*thisVal = NullUUIDFromString(&s.String)
	return nil
}

// MarshalJSON renders the value as a JSON string in canonical UUID
// form, or null when empty.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullUUID) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJSON, nil
	}
	return json.Marshal(thisVal.Val)
}

// UnmarshalJSON parses a JSON string containing a UUID, or the token
// null. Any other input is treated as a parse error.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullUUID) UnmarshalJSON(data []byte) error {
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
