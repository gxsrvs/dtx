package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

// NullOffsetDateTime is the nullable counterpart of OffsetDateTime. It
// shares the same serializer (formatOffsetDateTime / parseOffsetDateTime).
type NullOffsetDateTime struct {
	Val   time.Time
	Valid bool // Valid is true if Val is not NULL
}

func NewNullOffsetDateTime(value time.Time) NullOffsetDateTime {
	return NullOffsetDateTime{Val: value, Valid: true}
}

func NewNullOffsetDateTimeEmpty() NullOffsetDateTime {
	return NullOffsetDateTime{Valid: false}
}

func NullOffsetDateTimeFromString(strValue *string) NullOffsetDateTime {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullOffsetDateTimeEmpty()
	}
	parsed, err := parseOffsetDateTime(*strValue)
	if err != nil {
		return NewNullOffsetDateTimeEmpty()
	}
	return NewNullOffsetDateTime(parsed)
}

func (thisVal *NullOffsetDateTime) IsEmpty() bool {
	return !thisVal.Valid
}

func (thisVal NullOffsetDateTime) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatOffsetDateTime(thisVal.Val)
}

func (thisVal NullOffsetDateTime) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

func (thisVal *NullOffsetDateTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullOffsetDateTimeEmpty()
		return nil
	}
	*thisVal = NewNullOffsetDateTime(s.Time)
	return nil
}

func (thisVal NullOffsetDateTime) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(formatOffsetDateTime(thisVal.Val))
}

func (thisVal *NullOffsetDateTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := parseOffsetDateTime(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}

// Before reports whether the receiver is valid and before dt.
func (thisVal NullOffsetDateTime) Before(dt time.Time) bool {
	return thisVal.Valid && thisVal.Val.Before(dt)
}

// After reports whether the receiver is valid and after dt.
func (thisVal NullOffsetDateTime) After(dt time.Time) bool {
	return thisVal.Valid && thisVal.Val.After(dt)
}
