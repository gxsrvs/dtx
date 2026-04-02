package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	RuOnlyDateMask = "02.01.2006"
)

var dateFormats = []string{
	time.DateOnly,  // ISO
	RuOnlyDateMask, // RUS
}

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
	parsed, err := ParseDateFromString(*strValue)
	if err != nil {
		return NewNullDateEmpty()
	}
	return NewNullDate(*parsed)
}

func ParseDateFromString(strValue string) (*time.Time, error) {
	// Перебираем форматы пока не удастся успешно распарсить
	var parsed time.Time
	var err error
	for _, f := range dateFormats {
		parsed, err = time.ParseInLocation(f, strValue, time.Local)
		if err == nil {
			return &parsed, nil
		}
	}

	// Если ни один формат не подошёл
	return nil, errors.New("cannot parse date from " + strValue)
}

func DateToString(date time.Time) string {
	return date.Format(time.DateOnly)
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
	return thisVal.Val.Format(time.DateOnly)
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

	if reflect.TypeOf(value) == nil {
		*thisVal = NullDate{Val: s.Time}
	} else {
		*thisVal = NullDate{Val: s.Time, Valid: true}
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullDate) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(thisVal.Val.Format(time.DateOnly))
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDate) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := ParseDateFromString(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = *val
	return nil
}
