package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

// NullLocalDateTime is the nullable counterpart of LocalDateTime. It shares
// the same serializer (formatLocalDateTime / parseLocalDateTime).
type NullLocalDateTime struct {
	Val   LocalDateTime
	Valid bool // Valid is true if Val is not NULL
}

func NewNullLocalDateTime(value LocalDateTime) NullLocalDateTime {
	return NullLocalDateTime{Val: value, Valid: true}
}

func NewNullLocalDateTimeEmpty() NullLocalDateTime {
	return NullLocalDateTime{Valid: false}
}

func NullLocalDateTimeFromString(strValue *string) NullLocalDateTime {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullLocalDateTimeEmpty()
	}
	parsed, err := parseLocalDateTime(*strValue)
	if err != nil {
		return NewNullLocalDateTimeEmpty()
	}
	return NewNullLocalDateTime(parsed)
}

func (thisVal *NullLocalDateTime) IsEmpty() bool {
	return !thisVal.Valid
}

func (thisVal *NullLocalDateTime) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatLocalDateTime(thisVal.Val)
}

func (thisVal NullLocalDateTime) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val.ToTime(time.UTC), nil
}

func (thisVal *NullLocalDateTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullLocalDateTimeEmpty()
		return nil
	}
	*thisVal = NewNullLocalDateTime(localDateTimeFromTime(s.Time))
	return nil
}

func (thisVal NullLocalDateTime) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(formatLocalDateTime(thisVal.Val))
}

func (thisVal *NullLocalDateTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := parseLocalDateTime(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}
