package types

import (
	"database/sql/driver"
	"strings"
	"time"
)

// NullLocalTime is the nullable counterpart of LocalTime. It shares the
// same serializer (formatLocalTime / parseLocalTime). Valid reports
// whether Val holds a meaningful value (true) or a SQL/JSON NULL
// (false).
type NullLocalTime struct {
	Val   LocalTime
	Valid bool
}

// NewNullLocalTime constructs a valid NullLocalTime wrapping value. Use
// NullLocalTimeFromTime when starting from a time.Time.
func NewNullLocalTime(value LocalTime) NullLocalTime {
	return NullLocalTime{Val: value, Valid: true}
}

// NullLocalTimeFromTime constructs a valid NullLocalTime from a
// time.Time. The wall-clock fields are read in t's current location;
// the zone itself is dropped (LocalTime carries no zone by design).
func NullLocalTimeFromTime(value time.Time) NullLocalTime {
	return NullLocalTime{Val: LocalTimeFromTime(value), Valid: true}
}

// NewNullLocalTimeEmpty returns an invalid (NULL) NullLocalTime.
func NewNullLocalTimeEmpty() NullLocalTime {
	return NullLocalTime{Valid: false}
}

// NullLocalTimeFromString parses a local time from the string pointer.
// A nil pointer, an empty string, the tokens "null"/"nil"
// (case-insensitive), or a parse error all produce an invalid
// NullLocalTime.
func NullLocalTimeFromString(strValue *string) NullLocalTime {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullLocalTimeEmpty()
	}
	parsed, err := parseLocalTime(*strValue)
	if err != nil {
		return NewNullLocalTimeEmpty()
	}
	return NewNullLocalTime(parsed)
}

// IsEmpty reports whether the value is NULL (Valid == false).
func (thisVal *NullLocalTime) IsEmpty() bool {
	return !thisVal.Valid
}

// IsZero reports whether the value is NULL (Valid == false). Mirroring
// time.Time.IsZero, this also enables encoding/json's `omitzero` tag
// (Go 1.24+) to elide invalid wrappers from marshalled output.
func (thisVal NullLocalTime) IsZero() bool {
	return !thisVal.Valid
}

// ToString renders the value in the library's canonical local time
// form, or "" when NULL.
func (thisVal NullLocalTime) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatLocalTime(thisVal.Val)
}

// Value implements the database/sql/driver.Valuer interface. A NULL
// value is emitted as (nil, nil); the valid case is materialised via
// LocalTime.ToTime (a same-day time.Time in UTC).
func (thisVal NullLocalTime) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return thisVal.Val.ToTime(), nil
}

// Scan implements the database/sql.Scanner interface. NULL values
// produce an invalid NullLocalTime; non-NULL values are delegated to
// LocalTime.Scan.
func (thisVal *NullLocalTime) Scan(value interface{}) error {
	if value == nil {
		*thisVal = NewNullLocalTimeEmpty()
		return nil
	}
	var lt LocalTime
	if err := lt.Scan(value); err != nil {
		return err
	}
	*thisVal = NewNullLocalTime(lt)
	return nil
}

// MarshalJSON renders the value as a JSON string in canonical local
// time form, or null when empty. The valid path appends directly into
// a single buffer via time.Time.AppendFormat, skipping the
// intermediate string and reflection dispatch that
// json.Marshal(formatLocalTime(...)) would perform.
func (thisVal NullLocalTime) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJSON, nil
	}
	v := thisVal.Val
	t := time.Date(0, 1, 1, v.Hour, v.Minute, v.Second, v.Nanosec, time.UTC)
	// "15:04:05.999999999" is at most 18 bytes plus quotes.
	buf := make([]byte, 0, 24)
	buf = append(buf, '"')
	buf = t.AppendFormat(buf, "15:04:05.999999999")
	buf = append(buf, '"')
	return buf, nil
}

// UnmarshalJSON parses a JSON string containing a local time, or the
// token null. Any other input is treated as a parse error.
func (thisVal *NullLocalTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := parseLocalTime(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = val
	return nil
}

// Before reports whether thisVal precedes other under sortable NULL
// semantics: NULL is strictly less than any valid value, two NULLs
// compare equal (so neither is Before the other).
func (thisVal NullLocalTime) Before(other NullLocalTime) bool {
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
func (thisVal NullLocalTime) After(other NullLocalTime) bool {
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
func (thisVal NullLocalTime) Equal(other NullLocalTime) bool {
	if !thisVal.Valid {
		return !other.Valid
	}
	if !other.Valid {
		return false
	}
	return thisVal.Val.Equal(other.Val)
}

// Compare returns -1, 0, or +1 as thisVal is before, equal to, or
// after other under sortable NULL semantics.
func (thisVal NullLocalTime) Compare(other NullLocalTime) int {
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

// Add returns thisVal shifted by d. NULL propagates. No modulo-24h
// enforcement — see LocalTime.Add.
func (thisVal NullLocalTime) Add(d time.Duration) NullLocalTime {
	if !thisVal.Valid {
		return thisVal
	}
	return NullLocalTime{Val: thisVal.Val.Add(d), Valid: true}
}

// SubOk returns (thisVal − other, true) when both operands are valid,
// and (0, false) otherwise.
func (thisVal NullLocalTime) SubOk(other NullLocalTime) (time.Duration, bool) {
	if !thisVal.Valid || !other.Valid {
		return 0, false
	}
	return thisVal.Val.Sub(other.Val), true
}

// Truncate returns thisVal rounded down to the nearest multiple of d
// since the zero time. NULL propagates.
func (thisVal NullLocalTime) Truncate(d time.Duration) NullLocalTime {
	if !thisVal.Valid {
		return thisVal
	}
	return NullLocalTime{Val: thisVal.Val.Truncate(d), Valid: true}
}
