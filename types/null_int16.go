package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type NullInt16 struct {
	Val   int16
	Valid bool // Valid is true if Val is not NULL
}

func NewNullInt16(value int16) NullInt16 {
	return NullInt16{Val: value, Valid: true}
}

func NewNullInt16Empty() NullInt16 {
	return NullInt16{Valid: false}
}

func NullInt16FromString(strValue *string) NullInt16 {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullInt16Empty()
	}
	var result, err = strconv.ParseInt(*strValue, 10, 16)
	if err != nil {
		return NewNullInt16Empty()
	}
	return NewNullInt16(int16(result))
}

//goland:noinspection GoUnusedExportedFunction
func NullInt16FromNullString(str NullString) NullInt16 {
	if !str.Valid || str.Val == "" {
		return NewNullInt16Empty()
	}
	var result, err = strconv.ParseInt(str.Val, 10, 16)
	if err != nil {
		return NewNullInt16Empty()
	}
	return NewNullInt16(int16(result))
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt16) IsEmpty() bool {
	return !thisVal.Valid
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt16) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return fmt.Sprintf("%d", thisVal.Val)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullInt16) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return int64(thisVal.Val), nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt16) Scan(value interface{}) error {
	var s sql.NullInt16
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullInt16Empty()
		return nil
	}
	*thisVal = NewNullInt16(s.Int16)
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullInt16) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(thisVal.Val)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt16) UnmarshalJSON(data []byte) error {
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
	thisVal.Val = int16(val)
	return nil
}
