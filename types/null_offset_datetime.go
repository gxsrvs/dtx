package types

import (
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"
)

// NullOffsetDateTime is the nullable counterpart of OffsetDateTime. It
// shares the same serializer (formatOffsetDateTime /
// parseOffsetDateTime). Valid reports whether Val holds a meaningful
// value (true) or a SQL/JSON NULL (false).
type NullOffsetDateTime struct {
	Val   OffsetDateTime
	Valid bool
}

// NewNullOffsetDateTime constructs a valid NullOffsetDateTime wrapping
// value. Use NullOffsetDateTimeFromTime when starting from a
// time.Time.
func NewNullOffsetDateTime(value OffsetDateTime) NullOffsetDateTime {
	return NullOffsetDateTime{Val: value, Valid: true}
}

// NullOffsetDateTimeFromTime constructs a valid NullOffsetDateTime
// from a time.Time. Equivalent to
// NewNullOffsetDateTime(OffsetDateTime(value)).
func NullOffsetDateTimeFromTime(value time.Time) NullOffsetDateTime {
	return NullOffsetDateTime{Val: OffsetDateTime(value), Valid: true}
}

// NewNullOffsetDateTimeEmpty returns an invalid (NULL) NullOffsetDateTime.
func NewNullOffsetDateTimeEmpty() NullOffsetDateTime {
	return NullOffsetDateTime{Valid: false}
}

// NullOffsetDateTimeFromString parses an offset datetime from the
// string pointer. A nil pointer, an empty string, the tokens
// "null"/"nil" (case-insensitive), or a parse error all produce an
// invalid NullOffsetDateTime.
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
	return NullOffsetDateTimeFromTime(parsed)
}

// IsEmpty reports whether the value is NULL (Valid == false).
func (thisVal *NullOffsetDateTime) IsEmpty() bool {
	return !thisVal.Valid
}

// IsZero reports whether the value is NULL (Valid == false). Mirroring
// time.Time.IsZero, this also enables encoding/json's `omitzero` tag
// (Go 1.24+) to elide invalid wrappers from marshalled output.
func (thisVal NullOffsetDateTime) IsZero() bool {
	return !thisVal.Valid
}

// ToString renders the value in the library's canonical offset
// datetime form, or "" when NULL.
func (thisVal NullOffsetDateTime) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatOffsetDateTime(time.Time(thisVal.Val))
}

// Value implements the database/sql/driver.Valuer interface. A NULL
// value is emitted as (nil, nil); the valid case is emitted as the
// underlying time.Time.
func (thisVal NullOffsetDateTime) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return time.Time(thisVal.Val), nil
}

// Scan implements the database/sql.Scanner interface, delegating to
// sql.NullTime so that the driver's NULL signalling is honoured.
func (thisVal *NullOffsetDateTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullOffsetDateTimeEmpty()
		return nil
	}
	*thisVal = NullOffsetDateTimeFromTime(s.Time)
	return nil
}

// MarshalJSON renders the value as a JSON string in canonical offset
// datetime form, or null when empty. The valid path appends directly
// into a single buffer via time.Time.AppendFormat, skipping the
// intermediate string and reflection dispatch that
// json.Marshal(formatOffsetDateTime(...)) would perform.
func (thisVal NullOffsetDateTime) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJSON, nil
	}
	// "2006-01-02T15:04:05.999999999Z07:00" is at most 35 bytes plus quotes.
	buf := make([]byte, 0, 40)
	buf = append(buf, '"')
	buf = time.Time(thisVal.Val).AppendFormat(buf, offsetDateTimeFormat)
	buf = append(buf, '"')
	return buf, nil
}

// UnmarshalJSON parses a JSON string containing an offset datetime, or
// the token null. Any other input is treated as a parse error.
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
	thisVal.Val = OffsetDateTime(val)
	return nil
}

// Before reports whether thisVal precedes other under sortable NULL
// semantics: NULL is strictly less than any valid value, two NULLs
// compare equal (so neither is Before the other).
func (thisVal NullOffsetDateTime) Before(other NullOffsetDateTime) bool {
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
func (thisVal NullOffsetDateTime) After(other NullOffsetDateTime) bool {
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
func (thisVal NullOffsetDateTime) Equal(other NullOffsetDateTime) bool {
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
func (thisVal NullOffsetDateTime) Compare(other NullOffsetDateTime) int {
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

// Add returns thisVal shifted by d. NULL propagates.
func (thisVal NullOffsetDateTime) Add(d time.Duration) NullOffsetDateTime {
	if !thisVal.Valid {
		return thisVal
	}
	return NullOffsetDateTime{Val: thisVal.Val.Add(d), Valid: true}
}

// AddDate returns thisVal with the given years/months/days added.
// NULL propagates.
func (thisVal NullOffsetDateTime) AddDate(years, months, days int) NullOffsetDateTime {
	if !thisVal.Valid {
		return thisVal
	}
	return NullOffsetDateTime{Val: thisVal.Val.AddDate(years, months, days), Valid: true}
}

// SubOk returns (thisVal − other, true) when both operands are valid,
// and (0, false) otherwise.
func (thisVal NullOffsetDateTime) SubOk(other NullOffsetDateTime) (time.Duration, bool) {
	if !thisVal.Valid || !other.Valid {
		return 0, false
	}
	return thisVal.Val.Sub(other.Val), true
}

// Truncate returns thisVal rounded down to the nearest multiple of d
// since the zero time. NULL propagates.
func (thisVal NullOffsetDateTime) Truncate(d time.Duration) NullOffsetDateTime {
	if !thisVal.Valid {
		return thisVal
	}
	return NullOffsetDateTime{Val: thisVal.Val.Truncate(d), Valid: true}
}

// UnixOk returns the underlying instant as a Unix timestamp (seconds
// since 1970-01-01 UTC) and ok=true if the receiver is valid. For
// NULL the result is (0, false) — distinguishing NULL from the Unix
// epoch, which the legacy Unix() method (removed in v0.3.0) could
// not.
func (thisVal NullOffsetDateTime) UnixOk() (int64, bool) {
	if !thisVal.Valid {
		return 0, false
	}
	return thisVal.Val.Unix(), true
}

// UnixMilliOk returns the millisecond-precision Unix timestamp and
// ok=true if the receiver is valid; (0, false) for NULL.
func (thisVal NullOffsetDateTime) UnixMilliOk() (int64, bool) {
	if !thisVal.Valid {
		return 0, false
	}
	return thisVal.Val.UnixMilli(), true
}

// UnixMicroOk returns the microsecond-precision Unix timestamp and
// ok=true if the receiver is valid; (0, false) for NULL.
func (thisVal NullOffsetDateTime) UnixMicroOk() (int64, bool) {
	if !thisVal.Valid {
		return 0, false
	}
	return thisVal.Val.UnixMicro(), true
}

// UnixNanoOk returns the nanosecond-precision Unix timestamp and
// ok=true if the receiver is valid; (0, false) for NULL.
func (thisVal NullOffsetDateTime) UnixNanoOk() (int64, bool) {
	if !thisVal.Valid {
		return 0, false
	}
	return thisVal.Val.UnixNano(), true
}

// In returns the same instant rendered in loc, mirroring
// OffsetDateTime.In. NULL propagates: an invalid receiver is returned
// unchanged.
func (thisVal NullOffsetDateTime) In(loc *time.Location) NullOffsetDateTime {
	if !thisVal.Valid {
		return thisVal
	}
	return NullOffsetDateTime{Val: thisVal.Val.In(loc), Valid: true}
}

// UTC is shorthand for In(time.UTC). NULL propagates.
func (thisVal NullOffsetDateTime) UTC() NullOffsetDateTime {
	return thisVal.In(time.UTC)
}
