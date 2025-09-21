package types

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"
)

type DateTime time.Time

func NewDateTime(value time.Time) DateTime {
	return DateTime(value)
}

//goland:noinspection GoUnusedExportedFunction
func DateTimeFromString(strValue *string) (DateTime, error) {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return DateTime{}, errors.New("DateTime cannot be parsed from null value")
	}
	parsed, err := ParseDateTimeFromString(*strValue)
	if err != nil {
		return DateTime{}, errors.New("DateTime cannot be parsed from null value")
	}
	return NewDateTime(*parsed), nil
}

func ParseDateTimeFromString(strValue string) (*time.Time, error) {
	if len(strValue) < 11 {
		return nil, errors.New("cannot parse time from " + strValue)
	}
	sep := strValue[10]
	if sep != ' ' && sep != 'T' {
		return nil, errors.New("cannot parse time from " + strValue)
	}

	datePart := strValue[:10]
	parsedDate, err := ParseDateFromString(datePart)
	if err != nil {
		return nil, err
	}

	timePart := strValue[11:]
	parsedTime, err := ParseTimeFromString(timePart)
	if err != nil {
		return nil, err
	}

	// Теперь собираем итоговый time.Time
	result := AssembleDateTime(parsedDate, parsedTime, nil)
	return result, nil
}

func DateTimeToString(dateTime time.Time) string {
	return dateTime.Format(time.RFC3339)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *DateTime) ToString() string {
	return DateTimeToString(time.Time(*thisVal))
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *DateTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*thisVal = DateTime(s.Time)
	} else {
		*thisVal = DateTime(s.Time)
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal DateTime) MarshalJSON() ([]byte, error) {
	timeVal := time.Time(thisVal)
	return json.Marshal(timeVal.Format(time.RFC3339))
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *DateTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := ParseDateTimeFromString(s)
	if err != nil {
		return err
	}
	*thisVal = DateTime(*val)
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *DateTime) Before(dt time.Time) bool {
	return time.Time(*thisVal).Before(dt)
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *DateTime) After(dt time.Time) bool {
	return time.Time(*thisVal).After(dt)
}
