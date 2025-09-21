package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"strings"
)

type NullString struct {
	Val   string
	Valid bool // Valid is true if Val is not NULL
}

//goland:noinspection GoUnusedExportedFunction
func NewNullString(value string) NullString {
	return NullString{Val: value, Valid: true}
}

//goland:noinspection GoUnusedExportedFunction
func NewNullStringEmpty() NullString {
	return NullString{Valid: false}
}

//goland:noinspection GoMixedReceiverTypes,GoUnusedExportedFunction
func NSFromString(s string) sql.NullString {
	valid := true
	if s == "" {
		valid = false
	}
	return sql.NullString{
		String: s,
		Valid:  valid,
	}
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullString) IsEmpty() bool {
	return !thisVal.Valid
}

//goland:noinspection GoMixedReceiverTypes,GoUnusedExportedFunction
func GetNullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullString) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return thisVal.Val
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullString) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*thisVal = NullString{Val: s.String}
	} else {
		*thisVal = NullString{Val: s.String, Valid: true}
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullString) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(thisVal.Val)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullString) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		thisVal.Val = ""
		return nil
	}
	s := strings.Trim(sd, "\"")
	thisVal.Valid = true
	thisVal.Val = s
	return nil
}
