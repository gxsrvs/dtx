package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// Date is a not-null calendar date (no time component, no timezone).
//
// Date shares its serializer with NullDate: both render as "2006-01-02"
// and accept the same parser inputs (ISO 8601 and dd.MM.yyyy). Use Date
// for required JSON or SQL fields where NULL is not permitted; use
// NullDate for optional ones.
type Date time.Time

// NewDate wraps t as a Date. The time-of-day component is preserved on the
// underlying value but is dropped by the serializer (formatDate writes only
// "YYYY-MM-DD").
func NewDate(t time.Time) Date {
	return Date(t)
}

// DateFromString parses a calendar date using the formats supported by the
// library (ISO 8601 and dd.MM.yyyy).
func DateFromString(strValue string) (Date, error) {
	if strValue == "" {
		return Date{}, errors.New("cannot parse Date from empty string")
	}
	parsed, err := parseDate(strValue)
	if err != nil {
		return Date{}, err
	}
	return Date(parsed), nil
}

// AsTime returns the underlying time.Time.
func (thisVal Date) AsTime() time.Time {
	return time.Time(thisVal)
}

// ToString renders the date in the library's canonical format
// ("2006-01-02").
func (thisVal Date) ToString() string {
	return formatDate(time.Time(thisVal))
}

// Value implements the database/sql/driver.Valuer interface.
func (thisVal Date) Value() (driver.Value, error) {
	return time.Time(thisVal), nil
}

// Scan implements the database/sql.Scanner interface. A NULL value is
// rejected — Date cannot be empty.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *Date) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		return errors.New("cannot scan NULL into Date")
	}
	*thisVal = Date(s.Time)
	return nil
}

// MarshalJSON renders the date as a JSON string in the library's canonical
// format ("2006-01-02").
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(formatDate(time.Time(thisVal)))
}

// UnmarshalJSON parses a JSON string in any of the formats accepted by
// DateFromString. JSON null is rejected — Date cannot be empty.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *Date) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" {
		return errors.New("cannot unmarshal null into Date")
	}
	s := strings.Trim(sd, "\"")
	if s == "" {
		return errors.New("cannot unmarshal empty string into Date")
	}
	parsed, err := parseDate(s)
	if err != nil {
		return err
	}
	*thisVal = Date(parsed)
	return nil
}
