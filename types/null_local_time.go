package types

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

// NullLocalTime is the nullable counterpart of LocalTime. It shares the same
// serializer (formatLocalTime / parseLocalTime).
type NullLocalTime struct {
	Val   LocalTime
	Valid bool // Valid is true if Val is not NULL
}

func NewNullLocalTime(value LocalTime) NullLocalTime {
	return NullLocalTime{Val: value, Valid: true}
}

func NewNullLocalTimeEmpty() NullLocalTime {
	return NullLocalTime{Valid: false}
}

func NullLocalTimeFromString(strValue *string) NullLocalTime {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullLocalTimeEmpty()
	}
	parsed, err := parseLocalTime(*strValue)
	if err != nil {
		return NewNullLocalTimeEmpty()
	}
	return NewNullLocalTime(parsed)
}

func (thisVal *NullLocalTime) IsEmpty() bool {
	return !thisVal.Valid
}

func (thisVal *NullLocalTime) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatLocalTime(thisVal.Val)
}

func (thisVal NullLocalTime) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val.ToTime(), nil
}

func (thisVal *NullLocalTime) Scan(value interface{}) error {
	if value == nil {
		*thisVal = NewNullLocalTimeEmpty()
		return nil
	}
	var lt LocalTime
	if err := lt.Scan(value); err != nil {
		return err
	}
	*thisVal = NewNullLocalTime(lt)
	return nil
}

func (thisVal NullLocalTime) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(formatLocalTime(thisVal.Val))
}

func (thisVal *NullLocalTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := parseLocalTime(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}
