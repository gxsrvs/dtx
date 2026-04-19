package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestLocalTimeFromString(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want LocalTime
	}{
		{"hh:mm", "06:07", LocalTime{6, 7, 0, 0}},
		{"hh:mm:ss", "02:56:15", LocalTime{2, 56, 15, 0}},
		{"with millis", "02:56:15.123", LocalTime{2, 56, 15, 123_000_000}},
		{"with nanos", "02:56:15.123456789", LocalTime{2, 56, 15, 123_456_789}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := LocalTimeFromString(c.in)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got != c.want {
				t.Errorf("Expected %#v, got %#v", c.want, got)
			}
		})
	}
}

func TestLocalTimeFromStringRejectsTZ(t *testing.T) {
	rejected := []string{
		"02:56:15Z",
		"02:56:15+03:00",
		"02:56:15-05:00",
		"02:56:15+0300",
	}
	for _, in := range rejected {
		if _, err := LocalTimeFromString(in); err == nil {
			t.Errorf("Expected error for %q, got nil", in)
		}
	}
}

func TestLocalTimeFromStringErrors(t *testing.T) {
	if _, err := LocalTimeFromString(""); err == nil {
		t.Error("Expected error from empty string, got nil")
	}
	if _, err := LocalTimeFromString("not a time"); err == nil {
		t.Error("Expected error from invalid string, got nil")
	}
}

func TestLocalTimeToString(t *testing.T) {
	cases := []struct {
		in   LocalTime
		want string
	}{
		{LocalTime{2, 56, 15, 0}, "02:56:15"},
		{LocalTime{6, 7, 0, 0}, "06:07:00"},
		{LocalTime{2, 56, 15, 500_000_000}, "02:56:15.5"},
		{LocalTime{2, 56, 15, 123_000_000}, "02:56:15.123"},
	}
	for _, c := range cases {
		if got := c.in.ToString(); got != c.want {
			t.Errorf("Expected %s, got %s", c.want, got)
		}
	}
}

func TestLocalTimeRoundTripJSON(t *testing.T) {
	src := NewLocalTime(2, 56, 15, 123_000_000)
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"02:56:15.123"` {
		t.Errorf("Unexpected JSON: %s", string(data))
	}
	var dst LocalTime
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if src != dst {
		t.Errorf("Round-trip mismatch: %#v != %#v", src, dst)
	}
}

func TestLocalTimeUnmarshalRejectsNull(t *testing.T) {
	var v LocalTime
	if err := json.Unmarshal([]byte(`null`), &v); err == nil {
		t.Error("Expected error from JSON null, got nil")
	}
	if err := json.Unmarshal([]byte(`""`), &v); err == nil {
		t.Error("Expected error from empty JSON string, got nil")
	}
}

func TestLocalTimeFromTime(t *testing.T) {
	loc := time.FixedZone("UTC+03:00", 3*3600)
	src := time.Date(1969, time.July, 21, 2, 56, 15, 123_000_000, loc)
	got := LocalTimeFromTime(src)
	want := LocalTime{2, 56, 15, 123_000_000}
	if got != want {
		t.Errorf("Expected %#v, got %#v", want, got)
	}
}

func TestLocalTimeValueScan(t *testing.T) {
	src := NewLocalTime(2, 56, 15, 0)
	v, err := src.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	tv, ok := v.(time.Time)
	if !ok {
		t.Fatalf("Expected time.Time, got %T", v)
	}
	var dst LocalTime
	if err := dst.Scan(tv); err != nil {
		t.Fatalf("Scan(time.Time): %v", err)
	}
	if src != dst {
		t.Errorf("Round-trip mismatch: %#v != %#v", src, dst)
	}

	// Scan from string
	var fromStr LocalTime
	if err := fromStr.Scan("02:56:15"); err != nil {
		t.Fatalf("Scan(string): %v", err)
	}
	if fromStr != src {
		t.Errorf("Expected %#v, got %#v", src, fromStr)
	}

	// Scan from []byte
	var fromBytes LocalTime
	if err := fromBytes.Scan([]byte("02:56:15")); err != nil {
		t.Fatalf("Scan([]byte): %v", err)
	}
	if fromBytes != src {
		t.Errorf("Expected %#v, got %#v", src, fromBytes)
	}

	// Scan NULL → error
	var nullDst LocalTime
	if err := nullDst.Scan(nil); err == nil {
		t.Error("Expected error from NULL scan, got nil")
	}

	// Scan unsupported type → error
	var bogusDst LocalTime
	if err := bogusDst.Scan(42); err == nil {
		t.Error("Expected error from unsupported type, got nil")
	}
}
