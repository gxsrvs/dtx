package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

// NullOffsetTime is the nullable counterpart of OffsetTime. It shares the
// same serializer (formatOffsetTime / parseOffsetTime).
type NullOffsetTime struct {
	Val   time.Time
	Valid bool // Valid is true if Val is not NULL
}

func NewNullOffsetTime(value time.Time) NullOffsetTime {
	return NullOffsetTime{Val: value, Valid: true}
}

func NewNullOffsetTimeEmpty() NullOffsetTime {
	return NullOffsetTime{Valid: false}
}

func NullOffsetTimeFromString(strValue *string) NullOffsetTime {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullOffsetTimeEmpty()
	}
	parsed, err := parseOffsetTime(*strValue)
	if err != nil {
		return NewNullOffsetTimeEmpty()
	}
	return NewNullOffsetTime(parsed)
}

func (thisVal *NullOffsetTime) IsEmpty() bool {
	return !thisVal.Valid
}

func (thisVal NullOffsetTime) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatOffsetTime(thisVal.Val)
}

func (thisVal NullOffsetTime) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

func (thisVal *NullOffsetTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullOffsetTimeEmpty()
		return nil
	}
	*thisVal = NewNullOffsetTime(s.Time)
	return nil
}

func (thisVal NullOffsetTime) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(formatOffsetTime(thisVal.Val))
}

func (thisVal *NullOffsetTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := parseOffsetTime(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}
