package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"strings"
	"time"
)

type NullTime struct {
	Val   time.Time
	Valid bool // Valid is true if Val is not NULL
}

func NewNullTime(value time.Time) NullTime {
	return NullTime{Val: value, Valid: true}
}

func NewNullTimeEmpty() NullTime {
	return NullTime{Valid: false}
}

func NullTimeFromString(strValue *string) NullTime {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullTimeEmpty()
	}
	parsed, err := ParseTimeFromString(*strValue)
	if err != nil {
		return NewNullTimeEmpty()
	}
	return NewNullTime(*parsed)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullTime) IsEmpty() bool {
	return !thisVal.Valid
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullTime) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return thisVal.Val.String()
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullTime) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*thisVal = NullTime{Val: s.Time}
	} else {
		*thisVal = NullTime{Val: s.Time, Valid: true}
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullTime) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(thisVal.Val.Format(time.TimeOnly))
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := ParseTimeFromString(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = *val
	return nil
}
