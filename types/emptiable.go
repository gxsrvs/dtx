package types

import (
	"database/sql"
	"database/sql/driver"
	
	
)

type Emptiable interface {
	IsEmpty() bool
}

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
	case NullString, NullFloat, NullInt64, NullInt32, NullInt16, NullBool,
		NullTime, NullDate, NullIsoDate, NullDateTime:
		return v.(Emptiable).IsEmpty()
	case *NullString, *NullFloat, *NullInt64, *NullInt32, *NullInt16, *NullBool,
		*NullTime, *NullDate, *NullIsoDate, *NullDateTime:
		return v.(Emptiable).IsEmpty()
	default:
		return false
	}
}
