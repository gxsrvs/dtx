package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"strings"
	"time"
)

type NullIsoDate struct {
	Val   time.Time
	Valid bool // Valid is true if Val is not NULL
}

func NewNullIsoDate(value time.Time) NullIsoDate {
	return NullIsoDate{Val: value, Valid: true}
}

func NewNullIsoDateEmpty() NullIsoDate {
	return NullIsoDate{Valid: false}
}

//goland:noinspection GoUnusedExportedFunction
func NullIsoDateFromString(strValue *string) NullIsoDate {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullIsoDateEmpty()
	}
	var result, err = time.Parse(time.DateOnly, *strValue)
	if err != nil {
		return NewNullIsoDateEmpty()
	}
	return NewNullIsoDate(result)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullIsoDate) IsEmpty() bool {
	return !thisVal.Valid
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullIsoDate) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return thisVal.Val.Format(time.DateOnly)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullIsoDate) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullIsoDate) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*thisVal = NullIsoDate{Val: s.Time}
	} else {
		*thisVal = NullIsoDate{Val: s.Time, Valid: true}
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullIsoDate) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(thisVal.Val.Format(time.DateOnly))
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullIsoDate) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}
