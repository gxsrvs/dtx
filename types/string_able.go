package types

import (
	"fmt"
	
	
)

type ToStringAble interface {
	ToString() string
}

func ToString(val interface{}) string {
	//return fmt.Sprintf("%T", val)
	switch v := val.(type) {
	case NullString, NullFloat, NullInt64, NullInt32, NullInt16, NullBool, NullTime:
		return v.(ToStringAble).ToString()
	case *NullString, *NullFloat, *NullInt64, *NullInt32, *NullInt16, *NullBool, *NullTime:
		return v.(ToStringAble).ToString()
	default:
		return fmt.Sprintf("%v", v)
	}
}
