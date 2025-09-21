package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"strings"
	"time"
)

type NullDateTime struct {
	Val   time.Time
	Valid bool // Valid is true if Val is not NULL
}

func NewNullDateTime(value time.Time) NullDateTime {
	return NullDateTime{Val: value, Valid: true}
}

func NewNullDateTimeEmpty() NullDateTime {
	return NullDateTime{Valid: false}
}

//goland:noinspection GoUnusedExportedFunction
func NullDateTimeFromString(strValue *string) NullDateTime {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullDateTimeEmpty()
	}
	parsed, err := ParseDateTimeFromString(*strValue)
	if err != nil {
		return NewNullDateTimeEmpty()
	}
	return NewNullDateTime(*parsed)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDateTime) IsEmpty() bool {
	return !thisVal.Valid
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDateTime) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return DateTimeToString(thisVal.Val)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullDateTime) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDateTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*thisVal = NullDateTime{Val: s.Time}
	} else {
		*thisVal = NullDateTime{Val: s.Time, Valid: true}
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullDateTime) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(thisVal.Val.Format(time.RFC3339))
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDateTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := ParseDateTimeFromString(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = *val
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDateTime) Before(dt time.Time) bool {
	return thisVal.Valid && thisVal.Val.Before(dt)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDateTime) After(dt time.Time) bool {
	return thisVal.Valid && thisVal.Val.After(dt)
}
