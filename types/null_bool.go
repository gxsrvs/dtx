package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type NullBool struct {
	Val   bool
	Valid bool // Valid is true if Val is not NULL
}

func NewNullBool(value bool) NullBool {
	return NullBool{Val: value, Valid: true}
}

func NewNullBoolEmpty() NullBool {
	return NullBool{Valid: false}
}

//goland:noinspection GoUnusedExportedFunction
func NullBoolFromString(strValue *string) NullBool {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullBoolEmpty()
	}
	var result, err = strconv.ParseBool(*strValue)
	if err != nil {
		return NewNullBoolEmpty()
	}
	return NewNullBool(result)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullBool) IsEmpty() bool {
	return !thisVal.Valid
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullBool) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return fmt.Sprintf("%t", thisVal.Val)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullBool) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullBool) Scan(value interface{}) error {
	var s sql.NullBool
	if err := s.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*thisVal = NullBool{Val: s.Bool}
	} else {
		*thisVal = NullBool{Val: s.Bool, Valid: true}
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullBool) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(thisVal.Val)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullBool) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		thisVal.Val = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}
