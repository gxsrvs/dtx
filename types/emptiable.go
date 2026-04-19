package types

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
)

// Emptiable is the interface implemented by every nullable type in the
// package. IsEmpty reports whether the value is missing (NULL).
type Emptiable interface {
	IsEmpty() bool
}

// IsEmpty reports whether v should be treated as empty/NULL. It accepts
// nil, primitives, the standard library's sql.Null* family, and any value
// that implements Emptiable (including pointer-receiver implementations
// invoked on a value).
func IsEmpty(t interface{}) bool {
	if t == nil {
		return true
	}
	switch v := t.(type) {
	case string:
		return v == ""
	case int, int8, int16, int32, int64:
		return false
	case sql.NullString, sql.NullFloat64, sql.NullInt64, sql.NullInt32, sql.NullInt16, sql.NullBool, sql.NullTime:
		val, err := v.(driver.Valuer).Value()
		return val == nil || err != nil
	}
	if e, ok := t.(Emptiable); ok {
		return e.IsEmpty()
	}
	// Pointer-receiver Emptiable implementations are not reachable on a
	// non-pointer interface value. Wrap the value in a fresh pointer and
	// retry the assertion.
	rv := reflect.ValueOf(t)
	if rv.Kind() != reflect.Pointer {
		p := reflect.New(rv.Type())
		p.Elem().Set(rv)
		if e, ok := p.Interface().(Emptiable); ok {
			return e.IsEmpty()
		}
	}
	return false
}
