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

type NullFloat struct {
	Val   float64
	Valid bool // Valid is true if Val is not NULL
}

func NewNullFloat(value float64) NullFloat {
	return NullFloat{Val: value, Valid: true}
}

func NewNullFloatEmpty() NullFloat {
	return NullFloat{Valid: false}
}

//goland:noinspection GoUnusedExportedFunction
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

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullFloat) IsEmpty() bool {
	return !thisVal.Valid
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullFloat) ToString() string {
	if !thisVal.Valid {
		return "null"
	}
	return fmt.Sprintf("%f", thisVal.Val)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullFloat) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullFloat) Scan(value interface{}) error {
	var s sql.NullFloat64
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*thisVal = NullFloat{Val: s.Float64}
	} else {
		*thisVal = NullFloat{Val: s.Float64, Valid: true}
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullFloat) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(thisVal.Val)
}

//goland:noinspection GoMixedReceiverTypes
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
