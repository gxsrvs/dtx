package types

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

type NullUuid struct {
	Val   uuid.UUID
	Valid bool // Valid is true if UUID is not NULL
}

func NewNullUuid(val uuid.UUID) NullUuid {
	return NullUuid{
		Val:   val,
		Valid: true,
	}
}

func NewNullUuidEmpty() NullUuid {
	return NullUuid{
		Val:   uuid.UUID{},
		Valid: false,
	}
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullUuid) IsEmpty() bool {
	return !thisVal.Valid
}

//goland:noinspection GoUnusedExportedFunction
func (thisVal *NullUuid) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return thisVal.Val.String()
}

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

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullUuid) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*thisVal = NewNullUuidEmpty()
	} else {
		*thisVal = NullUuidFromString(&s.String)
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullUuid) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(thisVal.Val)

}

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
