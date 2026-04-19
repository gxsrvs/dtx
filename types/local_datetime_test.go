package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestLocalDateTimeFromString(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want LocalDateTime
	}{
		{"T form, no fractional", "1961-04-12T06:07:00",
			LocalDateTime{1961, time.April, 12, 6, 7, 0, 0}},
		{"space form, no fractional", "1961-04-12 06:07:00",
			LocalDateTime{1961, time.April, 12, 6, 7, 0, 0}},
		{"with millis", "1961-04-12T06:07:00.123",
			LocalDateTime{1961, time.April, 12, 6, 7, 0, 123_000_000}},
		{"with nanos", "1961-04-12T06:07:00.123456789",
			LocalDateTime{1961, time.April, 12, 6, 7, 0, 123_456_789}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := LocalDateTimeFromString(c.in)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got != c.want {
				t.Errorf("Expected %#v, got %#v", c.want, got)
			}
		})
	}
}

func TestLocalDateTimeFromStringRejectsTZ(t *testing.T) {
	rejected := []string{
		"1961-04-12T06:07:00Z",
		"1961-04-12T06:07:00+03:00",
		"1961-04-12T06:07:00-05:00",
		"1961-04-12T06:07:00+0300",
	}
	for _, in := range rejected {
		if _, err := LocalDateTimeFromString(in); err == nil {
			t.Errorf("Expected error for %q, got nil", in)
		}
	}
}

func TestLocalDateTimeFromStringErrors(t *testing.T) {
	if _, err := LocalDateTimeFromString(""); err == nil {
		t.Error("Expected error from empty string, got nil")
	}
	if _, err := LocalDateTimeFromString("not a datetime"); err == nil {
		t.Error("Expected error from invalid string, got nil")
	}
}

func TestLocalDateTimeToString(t *testing.T) {
	cases := []struct {
		in   LocalDateTime
		want string
	}{
		{LocalDateTime{1961, time.April, 12, 6, 7, 0, 0}, "1961-04-12T06:07:00"},
		{LocalDateTime{1961, time.April, 12, 6, 7, 0, 500_000_000}, "1961-04-12T06:07:00.5"},
		{LocalDateTime{1961, time.April, 12, 6, 7, 0, 123_000_000}, "1961-04-12T06:07:00.123"},
	}
	for _, c := range cases {
		if got := c.in.ToString(); got != c.want {
			t.Errorf("Expected %s, got %s", c.want, got)
		}
	}
}

func TestLocalDateTimeRoundTripJSON(t *testing.T) {
	src := NewLocalDateTime(1969, time.July, 21, 2, 56, 15, 123_000_000)
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"1969-07-21T02:56:15.123"` {
		t.Errorf("Unexpected JSON: %s", string(data))
	}
	var dst LocalDateTime
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if src != dst {
		t.Errorf("Round-trip mismatch: %#v != %#v", src, dst)
	}
}

func TestLocalDateTimeUnmarshalRejectsNull(t *testing.T) {
	var d LocalDateTime
	if err := json.Unmarshal([]byte(`null`), &d); err == nil {
		t.Error("Expected error from JSON null, got nil")
	}
	if err := json.Unmarshal([]byte(`""`), &d); err == nil {
		t.Error("Expected error from empty JSON string, got nil")
	}
}

func TestLocalDateTimeToTime(t *testing.T) {
	ldt := NewLocalDateTime(1969, time.July, 21, 2, 56, 15, 0)
	loc := time.FixedZone("UTC+03:00", 3*3600)
	got := ldt.ToTime(loc)
	want := time.Date(1969, time.July, 21, 2, 56, 15, 0, loc)
	if !got.Equal(want) {
		t.Errorf("Expected %v, got %v", want, got)
	}

	// nil location → UTC
	gotUTC := ldt.ToTime(nil)
	if gotUTC.Location() != time.UTC {
		t.Errorf("Expected UTC for nil loc, got %v", gotUTC.Location())
	}
}

func TestLocalDateTimeFromTime(t *testing.T) {
	loc := time.FixedZone("UTC+03:00", 3*3600)
	src := time.Date(1969, time.July, 21, 2, 56, 15, 123_000_000, loc)
	got := LocalDateTimeFromTime(src)
	want := LocalDateTime{1969, time.July, 21, 2, 56, 15, 123_000_000}
	if got != want {
		t.Errorf("Expected %#v, got %#v", want, got)
	}
}

func TestLocalDateTimeValueScan(t *testing.T) {
	src := NewLocalDateTime(1969, time.July, 21, 2, 56, 15, 0)
	v, err := src.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	tv, ok := v.(time.Time)
	if !ok {
		t.Fatalf("Expected time.Time, got %T", v)
	}
	var dst LocalDateTime
	if err := dst.Scan(tv); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if src != dst {
		t.Errorf("Round-trip mismatch: %#v != %#v", src, dst)
	}

	var nullDst LocalDateTime
	if err := nullDst.Scan(nil); err == nil {
		t.Error("Expected error from NULL scan, got nil")
	}
}
