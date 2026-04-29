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
//
// In(loc) and UTC() are not provided by design: Date is a calendar
// concept without a time-of-day, so zone reinterpretation is
// meaningless. Use OffsetDateTime if a zone-aware instant is required.
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
func (thisVal Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(formatDate(time.Time(thisVal)))
}

// UnmarshalJSON parses a JSON string in any of the formats accepted by
// DateFromString. JSON null is rejected — Date cannot be empty.
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

// Before reports whether thisVal precedes other.
func (thisVal Date) Before(other Date) bool {
	return time.Time(thisVal).Before(time.Time(other))
}

// After reports whether thisVal is after other.
func (thisVal Date) After(other Date) bool {
	return time.Time(thisVal).After(time.Time(other))
}

// Equal reports whether thisVal and other denote the same instant.
func (thisVal Date) Equal(other Date) bool {
	return time.Time(thisVal).Equal(time.Time(other))
}

// Compare returns -1, 0, or +1 as thisVal is before, equal to, or
// after other.
func (thisVal Date) Compare(other Date) int {
	return time.Time(thisVal).Compare(time.Time(other))
}

// Add returns the date shifted by d. Sub-day components of d are
// applied to the underlying time.Time but do not surface in the
// serialised form (Date renders only YYYY-MM-DD).
func (thisVal Date) Add(d time.Duration) Date {
	return Date(time.Time(thisVal).Add(d))
}

// AddDate returns the date with the given number of years, months,
// and days added.
func (thisVal Date) AddDate(years, months, days int) Date {
	return Date(time.Time(thisVal).AddDate(years, months, days))
}

// Sub returns the duration thisVal − other. Convert to days via
// `dur / (24 * time.Hour)` if needed.
func (thisVal Date) Sub(other Date) time.Duration {
	return time.Time(thisVal).Sub(time.Time(other))
}

// Unix returns the underlying instant as a Unix timestamp (seconds
// since 1970-01-01 UTC). For Date values constructed at midnight UTC
// this is the conventional "date as epoch seconds".
func (thisVal Date) Unix() int64 {
	return time.Time(thisVal).Unix()
}

// UnixMilli returns the underlying instant as milliseconds since
// 1970-01-01 UTC.
func (thisVal Date) UnixMilli() int64 {
	return time.Time(thisVal).UnixMilli()
}

// UnixMicro returns the underlying instant as microseconds since
// 1970-01-01 UTC.
func (thisVal Date) UnixMicro() int64 {
	return time.Time(thisVal).UnixMicro()
}

// UnixNano returns the underlying instant as nanoseconds since
// 1970-01-01 UTC.
func (thisVal Date) UnixNano() int64 {
	return time.Time(thisVal).UnixNano()
}

// Year returns the calendar year.
func (thisVal Date) Year() int { return time.Time(thisVal).Year() }

// Month returns the calendar month (1-based, time.January..December).
func (thisVal Date) Month() time.Month { return time.Time(thisVal).Month() }

// Day returns the day of the month (1-31).
func (thisVal Date) Day() int { return time.Time(thisVal).Day() }

// YearDay returns the day-of-year (1-365 / 1-366 in leap years).
func (thisVal Date) YearDay() int { return time.Time(thisVal).YearDay() }
