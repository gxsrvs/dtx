package types

import (
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"
)

// NullOffsetTime is the nullable counterpart of OffsetTime. It shares
// the same serializer (formatOffsetTime / parseOffsetTime). Valid
// reports whether Val holds a meaningful value (true) or a SQL/JSON
// NULL (false).
type NullOffsetTime struct {
	Val   OffsetTime
	Valid bool
}

// NewNullOffsetTime constructs a valid NullOffsetTime wrapping value.
// Use NullOffsetTimeFromTime when starting from a time.Time.
func NewNullOffsetTime(value OffsetTime) NullOffsetTime {
	return NullOffsetTime{Val: value, Valid: true}
}

// NullOffsetTimeFromTime constructs a valid NullOffsetTime from a
// time.Time. Equivalent to NewNullOffsetTime(OffsetTime(value)).
func NullOffsetTimeFromTime(value time.Time) NullOffsetTime {
	return NullOffsetTime{Val: OffsetTime(value), Valid: true}
}

// NewNullOffsetTimeEmpty returns an invalid (NULL) NullOffsetTime.
func NewNullOffsetTimeEmpty() NullOffsetTime {
	return NullOffsetTime{Valid: false}
}

// NullOffsetTimeFromString parses an offset time from the string
// pointer. A nil pointer, an empty string, the tokens "null"/"nil"
// (case-insensitive), or a parse error all produce an invalid
// NullOffsetTime.
func NullOffsetTimeFromString(strValue *string) NullOffsetTime {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullOffsetTimeEmpty()
	}
	parsed, err := parseOffsetTime(*strValue)
	if err != nil {
		return NewNullOffsetTimeEmpty()
	}
	return NullOffsetTimeFromTime(parsed)
}

// IsEmpty reports whether the value is NULL (Valid == false).
func (thisVal *NullOffsetTime) IsEmpty() bool {
	return !thisVal.Valid
}

// IsZero reports whether the value is NULL (Valid == false). Mirroring
// time.Time.IsZero, this also enables encoding/json's `omitzero` tag
// (Go 1.24+) to elide invalid wrappers from marshalled output.
func (thisVal NullOffsetTime) IsZero() bool {
	return !thisVal.Valid
}

// ToString renders the value in the library's canonical offset time
// form, or "" when NULL.
func (thisVal NullOffsetTime) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return formatOffsetTime(time.Time(thisVal.Val))
}

// Value implements the database/sql/driver.Valuer interface. A NULL
// value is emitted as (nil, nil); the valid case is emitted as the
// underlying time.Time.
func (thisVal NullOffsetTime) Value() (driver.Value, error) {
	if !thisVal.Valid {
		return nil, nil
	}
	return time.Time(thisVal.Val), nil
}

// Scan implements the database/sql.Scanner interface, delegating to
// sql.NullTime so that the driver's NULL signalling is honoured.
func (thisVal *NullOffsetTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}
	if !s.Valid {
		*thisVal = NewNullOffsetTimeEmpty()
		return nil
	}
	*thisVal = NullOffsetTimeFromTime(s.Time)
	return nil
}

// MarshalJSON renders the value as a JSON string in canonical offset
// time form, or null when empty. The valid path appends directly into
// a single buffer via time.Time.AppendFormat, skipping the
// intermediate string and reflection dispatch that
// json.Marshal(formatOffsetTime(...)) would perform.
func (thisVal NullOffsetTime) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJSON, nil
	}
	// "15:04:05.999999999Z07:00" is at most 24 bytes plus quotes.
	buf := make([]byte, 0, 32)
	buf = append(buf, '"')
	buf = time.Time(thisVal.Val).AppendFormat(buf, offsetTimeFormat)
	buf = append(buf, '"')
	return buf, nil
}

// UnmarshalJSON parses a JSON string containing an offset time, or the
// token null. Any other input is treated as a parse error.
func (thisVal *NullOffsetTime) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		return nil
	}
	s := strings.Trim(sd, "\"")
	val, err := parseOffsetTime(s)
	if err != nil {
		return err
	}
	thisVal.Valid = true
	thisVal.Val = OffsetTime(val)
	return nil
}

// In returns the same instant rendered in loc, mirroring OffsetTime.In.
// NULL propagates: an invalid receiver is returned unchanged.
func (thisVal NullOffsetTime) In(loc *time.Location) NullOffsetTime {
	if !thisVal.Valid {
		return thisVal
	}
	return NullOffsetTime{Val: thisVal.Val.In(loc), Valid: true}
}

// UTC is shorthand for In(time.UTC). NULL propagates.
func (thisVal NullOffsetTime) UTC() NullOffsetTime {
	return thisVal.In(time.UTC)
}

// Before reports whether thisVal precedes other under sortable NULL
// semantics: NULL is strictly less than any valid value, two NULLs
// compare equal (so neither is Before the other).
func (thisVal NullOffsetTime) Before(other NullOffsetTime) bool {
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
func (thisVal NullOffsetTime) After(other NullOffsetTime) bool {
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
func (thisVal NullOffsetTime) Equal(other NullOffsetTime) bool {
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
func (thisVal NullOffsetTime) Compare(other NullOffsetTime) int {
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
// enforcement — see OffsetTime.Add.
func (thisVal NullOffsetTime) Add(d time.Duration) NullOffsetTime {
	if !thisVal.Valid {
		return thisVal
	}
	return NullOffsetTime{Val: thisVal.Val.Add(d), Valid: true}
}

// SubOk returns (thisVal − other, true) when both operands are valid,
// and (0, false) otherwise.
func (thisVal NullOffsetTime) SubOk(other NullOffsetTime) (time.Duration, bool) {
	if !thisVal.Valid || !other.Valid {
		return 0, false
	}
	return thisVal.Val.Sub(other.Val), true
}

// Truncate returns thisVal rounded down to the nearest multiple of d
// since the zero time. NULL propagates.
func (thisVal NullOffsetTime) Truncate(d time.Duration) NullOffsetTime {
	if !thisVal.Valid {
		return thisVal
	}
	return NullOffsetTime{Val: thisVal.Val.Truncate(d), Valid: true}
}
