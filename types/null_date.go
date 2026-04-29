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
	Val   Date
	Valid bool
}

// NewNullDate constructs a valid NullDate wrapping value. Use
// NullDateFromTime when starting from a time.Time.
func NewNullDate(value Date) NullDate {
	return NullDate{Val: value, Valid: true}
}

// NullDateFromTime constructs a valid NullDate from a time.Time.
// Equivalent to NewNullDate(Date(value)).
func NullDateFromTime(value time.Time) NullDate {
	return NullDate{Val: Date(value), Valid: true}
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
	return NullDateFromTime(parsed)
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
func (thisVal *NullDate) IsEmpty() bool {
	return !thisVal.Valid
}

// IsZero reports whether the value is NULL (Valid == false). Mirroring
// time.Time.IsZero, this also enables encoding/json's `omitzero` tag
// (Go 1.24+) to elide invalid wrappers from marshalled output.
func (thisVal NullDate) IsZero() bool {
	return !thisVal.Valid
}

// ToString renders the date in ISO 8601 form, or "" when NULL.
func (thisVal NullDate) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatDate(time.Time(thisVal.Val))
}

// Value implements the database/sql/driver.Valuer interface. A NULL
// value is emitted as (nil, nil); the valid case is emitted as the
// underlying time.Time (drivers then format it as DATE).
func (thisVal NullDate) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return time.Time(thisVal.Val), nil
}

// Scan implements the database/sql.Scanner interface, delegating to
// sql.NullTime so that the driver's NULL signalling is honoured.
func (thisVal *NullDate) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullDateEmpty()
		return nil
	}
	*thisVal = NullDateFromTime(s.Time)
	return nil
}

// MarshalJSON renders the date as a JSON string in ISO 8601 form, or
// null when empty. The valid path writes directly to a single buffer
// via time.Time.AppendFormat, skipping the intermediate string and
// reflection dispatch that json.Marshal would perform.
func (thisVal NullDate) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJSON, nil
	}
	// "2006-01-02" is 10 bytes, plus the two surrounding quotes.
	buf := make([]byte, 0, len(time.DateOnly)+2)
	buf = append(buf, '"')
	buf = time.Time(thisVal.Val).AppendFormat(buf, time.DateOnly)
	buf = append(buf, '"')
	return buf, nil
}

// UnmarshalJSON parses a JSON string using the accepted date formats
// (ISO 8601, dd.MM.yyyy) or the token null.
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
	thisVal.Val = Date(val)
	return nil
}

// Before reports whether thisVal precedes other under sortable NULL
// semantics: NULL is strictly less than any valid value, two NULLs
// compare equal (so neither is Before the other).
func (thisVal NullDate) Before(other NullDate) bool {
	if !thisVal.Valid {
		return other.Valid
	}
	if !other.Valid {
		return false
	}
	return thisVal.Val.Before(other.Val)
}

// After reports whether thisVal is after other under sortable NULL
// semantics: NULL is never after any value (including another NULL).
func (thisVal NullDate) After(other NullDate) bool {
	if !thisVal.Valid {
		return false
	}
	if !other.Valid {
		return true
	}
	return thisVal.Val.After(other.Val)
}

// Equal reports whether thisVal and other are equal under sortable
// NULL semantics: two NULLs are equal; NULL is never equal to a
// valid value.
func (thisVal NullDate) Equal(other NullDate) bool {
	if !thisVal.Valid {
		return !other.Valid
	}
	if !other.Valid {
		return false
	}
	return thisVal.Val.Equal(other.Val)
}

// Compare returns -1, 0, or +1 as thisVal is before, equal to, or
// after other under sortable NULL semantics: NULL precedes any valid
// value, two NULLs compare equal.
func (thisVal NullDate) Compare(other NullDate) int {
	if !thisVal.Valid {
		if !other.Valid {
			return 0
		}
		return -1
	}
	if !other.Valid {
		return +1
	}
	return thisVal.Val.Compare(other.Val)
}

// Add returns the date shifted by d. NULL propagates: an invalid
// receiver is returned unchanged.
func (thisVal NullDate) Add(d time.Duration) NullDate {
	if !thisVal.Valid {
		return thisVal
	}
	return NullDate{Val: thisVal.Val.Add(d), Valid: true}
}

// AddDate returns the date with the given years/months/days added.
// NULL propagates.
func (thisVal NullDate) AddDate(years, months, days int) NullDate {
	if !thisVal.Valid {
		return thisVal
	}
	return NullDate{Val: thisVal.Val.AddDate(years, months, days), Valid: true}
}

// SubOk returns (thisVal − other, true) when both operands are
// valid, and (0, false) otherwise.
func (thisVal NullDate) SubOk(other NullDate) (time.Duration, bool) {
	if !thisVal.Valid || !other.Valid {
		return 0, false
	}
	return thisVal.Val.Sub(other.Val), true
}

// UnixOk returns the underlying instant as a Unix timestamp (seconds
// since 1970-01-01 UTC) and ok=true if the receiver is valid;
// (0, false) for NULL.
func (thisVal NullDate) UnixOk() (int64, bool) {
	if !thisVal.Valid {
		return 0, false
	}
	return thisVal.Val.Unix(), true
}

// UnixMilliOk returns the millisecond-precision Unix timestamp and
// ok=true if the receiver is valid; (0, false) for NULL.
func (thisVal NullDate) UnixMilliOk() (int64, bool) {
	if !thisVal.Valid {
		return 0, false
	}
	return thisVal.Val.UnixMilli(), true
}

// UnixMicroOk returns the microsecond-precision Unix timestamp and
// ok=true if the receiver is valid; (0, false) for NULL.
func (thisVal NullDate) UnixMicroOk() (int64, bool) {
	if !thisVal.Valid {
		return 0, false
	}
	return thisVal.Val.UnixMicro(), true
}

// UnixNanoOk returns the nanosecond-precision Unix timestamp and
// ok=true if the receiver is valid; (0, false) for NULL.
func (thisVal NullDate) UnixNanoOk() (int64, bool) {
	if !thisVal.Valid {
		return 0, false
	}
	return thisVal.Val.UnixNano(), true
}
