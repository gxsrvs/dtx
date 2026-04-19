package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type NullInt64 struct {
	Val   int64
	Valid bool // Valid is true if Val is not NULL
}

func NewNullInt64(value int64) NullInt64 {
	return NullInt64{Val: value, Valid: true}
}

func NewNullInt64Empty() NullInt64 {
	return NullInt64{Valid: false}
}

func NullInt64FromString(strValue *string) NullInt64 {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullInt64Empty()
	}
	result, err := strconv.ParseInt(*strValue, 10, 64)
	if err != nil {
		return NewNullInt64Empty()
	}
	return NewNullInt64(result)
}

func NullInt64FromNullString(str NullString) NullInt64 {
	if !str.Valid || str.Val == "" {
		return NewNullInt64Empty()
	}
	var result, err = strconv.ParseInt(str.Val, 10, 64)
	if err != nil {
		return NewNullInt64Empty()
	}
	return NewNullInt64(int64(result))
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt64) IsEmpty() bool {
	return !thisVal.Valid
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt64) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return fmt.Sprintf("%d", thisVal.Val)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullInt64) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt64) Scan(value interface{}) error {
	var s sql.NullInt64
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullInt64Empty()
		return nil
	}
	*thisVal = NewNullInt64(s.Int64)
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullInt64) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(thisVal.Val)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullInt64) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		thisVal.Val = 0
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}
