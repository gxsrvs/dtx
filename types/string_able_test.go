package types

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestToString(t *testing.T) {
	d := time.Date(1969, 7, 20, 0, 0, 0, 0, time.UTC)
	odt := time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC)
	u, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")

	cases := []struct {
		name string
		in   interface{}
		want string
	}{
		{"NullString valid", NewNullString("hello"), "hello"},
		{"NullString empty", NewNullStringEmpty(), ""},
		{"NullInt64", NewNullInt64(42), "42"},
		{"NullInt32", NewNullInt32(42), "42"},
		{"NullInt16", NewNullInt16(42), "42"},
		{"NullBool", NewNullBool(true), "true"},
		{"NullDate", NewNullDate(d), "1969-07-20"},
		{"NullOffsetTime", NewNullOffsetTime(time.Date(0, 1, 1, 20, 17, 40, 0, time.UTC)), "20:17:40Z"},
		{"NullOffsetDateTime", NewNullOffsetDateTime(odt), "1969-07-20T20:17:40Z"},
		{"NullLocalDateTime", NewNullLocalDateTime(LocalDateTimeFromTime(odt)), "1969-07-20T20:17:40"},
		{"NullLocalTime", NewNullLocalTime(LocalTimeFromTime(odt)), "20:17:40"},
		{"NullUUID", NewNullUUID(u), "550e8400-e29b-41d4-a716-446655440000"},
		{"NullDecimal", NewNullDecimal(decimal.NewFromFloat(1.5)), "1.5"},

		{"*NullString", ptr(NewNullString("hello")), "hello"},
		{"*NullInt64", ptr(NewNullInt64(7)), "7"},

		{"plain string fallback", "raw", "raw"},
		{"plain int fallback", 123, "123"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ToString(c.in); got != c.want {
				t.Errorf("ToString(%v) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
