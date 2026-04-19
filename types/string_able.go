package types

import (
	"fmt"
	"reflect"
)

// ToStringAble is the interface implemented by every nullable type in the
// package. ToString renders the value in its canonical textual form (an
// empty string when the value is missing/NULL).
type ToStringAble interface {
	ToString() string
}

// ToString renders v using its ToString method when available (including
// pointer-receiver implementations invoked on a value), and falls back to
// fmt.Sprintf("%v", ...) otherwise.
func ToString(val interface{}) string {
	if val == nil {
		return ""
	}
	if s, ok := val.(ToStringAble); ok {
		return s.ToString()
	}
	rv := reflect.ValueOf(val)
	if rv.Kind() != reflect.Pointer {
		p := reflect.New(rv.Type())
		p.Elem().Set(rv)
		if s, ok := p.Interface().(ToStringAble); ok {
			return s.ToString()
		}
	}
	return fmt.Sprintf("%v", val)
}
