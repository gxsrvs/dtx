package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

type NullDate struct {
	Val   time.Time
	Valid bool // Valid is true if Val is not NULL
}

func NewNullDate(value time.Time) NullDate {
	return NullDate{Val: value, Valid: true}
}

func NewNullDateEmpty() NullDate {
	return NullDate{Valid: false}
}

func NullDateFromString(strValue *string) NullDate {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullDateEmpty()
	}
	parsed, err := parseDate(*strValue)
	if err != nil {
		return NewNullDateEmpty()
	}
	return NewNullDate(parsed)
}

// ParseDateFromString parses a calendar date from a string using the formats
// supported by the library (ISO 8601 and dd.MM.yyyy).
func ParseDateFromString(strValue string) (*time.Time, error) {
	parsed, err := parseDate(strValue)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// DateToString renders a calendar date in the library's canonical format
// ("2006-01-02").
func DateToString(date time.Time) string {
	return formatDate(date)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDate) IsEmpty() bool {
	return !thisVal.Valid
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDate) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatDate(thisVal.Val)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullDate) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDate) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullDateEmpty()
		return nil
	}
	*thisVal = NewNullDate(s.Time)
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullDate) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(formatDate(thisVal.Val))
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDate) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := parseDate(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}
