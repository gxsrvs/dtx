package types

import (
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"
)

// NullDate is a nullable calendar date (year/month/day with no
// time-of-day or zone component). Valid reports whether Val holds a
// meaningful value (true) or a SQL/JSON NULL (false).
type NullDate struct {
	Val   time.Time
	Valid bool
}

// NewNullDate constructs a valid NullDate wrapping value.
func NewNullDate(value time.Time) NullDate {
	return NullDate{Val: value, Valid: true}
}

// NewNullDateEmpty returns an invalid (NULL) NullDate.
func NewNullDateEmpty() NullDate {
	return NullDate{Valid: false}
}

// NullDateFromString parses a date from the string pointer using the
// library's accepted formats (ISO 8601 "2006-01-02" and the Russian
// "dd.MM.yyyy"). A nil pointer, an empty string, the tokens
// "null"/"nil" (case-insensitive), or a parse error all produce an
// invalid NullDate.
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

// ParseDateFromString parses a calendar date from a string using the
// formats supported by the library (ISO 8601 and dd.MM.yyyy). The
// returned pointer is never nil on success.
func ParseDateFromString(strValue string) (*time.Time, error) {
	parsed, err := parseDate(strValue)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// DateToString renders a calendar date in the library's canonical ISO
// 8601 form ("2006-01-02").
func DateToString(date time.Time) string {
	return formatDate(date)
}

// IsEmpty reports whether the value is NULL (Valid == false).
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDate) IsEmpty() bool {
	return !thisVal.Valid
}

// IsZero reports whether the value is NULL (Valid == false). Mirroring
// time.Time.IsZero, this also enables encoding/json's `omitzero` tag
// (Go 1.24+) to elide invalid wrappers from marshalled output.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullDate) IsZero() bool {
	return !thisVal.Valid
}

// ToString renders the date in ISO 8601 form, or "" when NULL.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDate) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatDate(thisVal.Val)
}

// Value implements the database/sql/driver.Valuer interface. A NULL
// value is emitted as (nil, nil); the valid case is emitted as the
// underlying time.Time (drivers then format it as DATE).
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullDate) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val, nil
}

// Scan implements the database/sql.Scanner interface, delegating to
// sql.NullTime so that the driver's NULL signalling is honoured.
//
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

// MarshalJSON renders the date as a JSON string in ISO 8601 form, or
// null when empty. The valid path writes directly to a single buffer
// via time.Time.AppendFormat, skipping the intermediate string and
// reflection dispatch that json.Marshal would perform.
//
//goland:noinspection GoMixedReceiverTypes
func (thisVal NullDate) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJSON, nil
	}
	// "2006-01-02" is 10 bytes, plus the two surrounding quotes.
	buf := make([]byte, 0, len(time.DateOnly)+2)
	buf = append(buf, '"')
	buf = thisVal.Val.AppendFormat(buf, time.DateOnly)
	buf = append(buf, '"')
	return buf, nil
}

// UnmarshalJSON parses a JSON string using the accepted date formats
// (ISO 8601, dd.MM.yyyy) or the token null.
//
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
