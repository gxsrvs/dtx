package types

import (
	"database/sql"
	"testing"
	"time"
)

func TestIsEmpty(t *testing.T) {
	cases := []struct {
		name string
		in   interface{}
		want bool
	}{
		{"nil interface", nil, true},
		{"empty string", "", true},
		{"non-empty string", "x", false},
		{"int zero", 0, false},
		{"int non-zero", 42, false},
		{"int8", int8(1), false},
		{"int16", int16(1), false},
		{"int32", int32(1), false},
		{"int64", int64(1), false},

		{"sql.NullString invalid", sql.NullString{Valid: false}, true},
		{"sql.NullString valid", sql.NullString{String: "x", Valid: true}, false},
		{"sql.NullFloat64 invalid", sql.NullFloat64{Valid: false}, true},
		{"sql.NullInt64 valid", sql.NullInt64{Int64: 1, Valid: true}, false},
		{"sql.NullBool invalid", sql.NullBool{Valid: false}, true},
		{"sql.NullTime invalid", sql.NullTime{Valid: false}, true},
		{"sql.NullByte invalid", sql.NullByte{Valid: false}, true},
		{"sql.NullByte valid", sql.NullByte{Byte: 0x42, Valid: true}, false},

		{"NullString empty", NewNullStringEmpty(), true},
		{"NullString valid", NewNullString("x"), false},
		{"NullFloat empty", NewNullFloatEmpty(), true},
		{"NullFloat valid", NewNullFloat(1.5), false},
		{"NullInt16 empty", NewNullInt16Empty(), true},
		{"NullInt32 empty", NewNullInt32Empty(), true},
		{"NullInt64 empty", NewNullInt64Empty(), true},
		{"NullBool empty", NewNullBoolEmpty(), true},
		{"NullDate empty", NewNullDateEmpty(), true},
		{"NullOffsetTime empty", NewNullOffsetTimeEmpty(), true},
		{"NullOffsetDateTime empty", NewNullOffsetDateTimeEmpty(), true},
		{"NullLocalDateTime empty", NewNullLocalDateTimeEmpty(), true},
		{"NullLocalTime empty", NewNullLocalTimeEmpty(), true},

		{"*NullString empty", ptr(NewNullStringEmpty()), true},
		{"*NullDate valid", ptr(NullDateFromTime(time.Now())), false},
		{"*NullOffsetTime empty", ptr(NewNullOffsetTimeEmpty()), true},

		{"unknown type fallback", time.Time{}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsEmpty(c.in); got != c.want {
				t.Errorf("IsEmpty(%v) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func ptr[T any](v T) *T { return &v }
