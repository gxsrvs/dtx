package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type NullInt32 struct {
	Val   int32
	Valid bool // Valid is true if Val is not NULL
}

func NewNullInt32(value int32) NullInt32 {
	return NullInt32{Val: value, Valid: true}
}

func NewNullInt32Empty() NullInt32 {
	return NullInt32{Valid: false}
}

func NullInt32FromString(strValue *string) NullInt32 {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullInt32Empty()
	}
	var result, err = strconv.ParseInt(*strValue, 10, 32)
	if err != nil {
		return NewNullInt32Empty()
	}
	return NewNullInt32(int32(result))
}

func NullInt32FromNullString(str NullString) NullInt32 {
	if !str.Valid || str.Val == "" {
		return NewNullInt32Empty()
	}
	var result, err = strconv.ParseInt(str.Val, 10, 32)
	if err != nil {
		return NewNullInt32Empty()
	}
	return NewNullInt32(int32(result))
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt32) IsEmpty() bool {
	return !thisVal.Valid
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt32) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return fmt.Sprintf("%d", thisVal.Val)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullInt32) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return int64(thisVal.Val), nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt32) Scan(value interface{}) error {
	var s sql.NullInt32
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullInt32Empty()
		return nil
	}
	*thisVal = NewNullInt32(s.Int32)
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullInt32) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(thisVal.Val)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt32) UnmarshalJSON(data []byte) error {
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
	thisVal.Val = int32(val)
	return nil
}
