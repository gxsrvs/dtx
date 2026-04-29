package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestOffsetDateTimeFromString(t *testing.T) {
	o, err := OffsetDateTimeFromString("1961-04-12T09:07:00+03:00")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	got := o.AsTime()
	if got.Year() != 1961 || got.Month() != 4 || got.Day() != 12 ||
		got.Hour() != 9 || got.Minute() != 7 || got.Second() != 0 {
		t.Errorf("Unexpected datetime value: %v", got)
	}
	_, off := got.Zone()
	if off != 3*3600 {
		t.Errorf("Expected offset +03:00, got %d", off)
	}

	if _, err := OffsetDateTimeFromString(""); err == nil {
		t.Error("Expected error from empty string, got nil")
	}
	if _, err := OffsetDateTimeFromString("not a datetime"); err == nil {
		t.Error("Expected error from invalid string, got nil")
	}
}

func TestOffsetDateTimeRoundTripJSON(t *testing.T) {
	src := NewOffsetDateTime(time.Date(1961, 4, 12, 9, 7, 0, 0, time.FixedZone("UTC+03:00", 3*3600)))
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"1961-04-12T09:07:00+03:00"` {
		t.Errorf("Unexpected JSON: %s", string(data))
	}

	var dst OffsetDateTime
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !dst.AsTime().Equal(src.AsTime()) {
		t.Errorf("Round-trip mismatch: %v != %v", dst.AsTime(), src.AsTime())
	}
}

func TestOffsetDateTimeMarshalUTC(t *testing.T) {
	v := NewOffsetDateTime(time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC))
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"1969-07-20T20:17:40Z"` {
		t.Errorf("Expected Z suffix for UTC, got %s", string(data))
	}
}

func TestOffsetDateTimeUnmarshalRejectsNull(t *testing.T) {
	var v OffsetDateTime
	if err := json.Unmarshal([]byte(`null`), &v); err == nil {
		t.Error("Expected error from JSON null, got nil")
	}
	if err := json.Unmarshal([]byte(`""`), &v); err == nil {
		t.Error("Expected error from empty JSON string, got nil")
	}
}

func TestOffsetDateTimeValueScan(t *testing.T) {
	src := time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC)
	o := NewOffsetDateTime(src)
	v, err := o.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	tv, ok := v.(time.Time)
	if !ok {
		t.Fatalf("Expected time.Time, got %T", v)
	}
	if !tv.Equal(src) {
		t.Errorf("Expected %v, got %v", src, tv)
	}

	var dst OffsetDateTime
	if err := dst.Scan(src); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if !dst.AsTime().Equal(src) {
		t.Errorf("Expected %v, got %v", src, dst.AsTime())
	}

	var nullDst OffsetDateTime
	if err := nullDst.Scan(nil); err == nil {
		t.Error("Expected error from NULL scan, got nil")
	}
}

func TestOffsetDateTime_ToString(t *testing.T) {
	v := NewOffsetDateTime(time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC))
	if got := v.ToString(); got != "1969-07-20T20:17:40Z" {
		t.Errorf("Expected '1969-07-20T20:17:40Z', got %q", got)
	}
}

// TestOffsetDateTime_Unix verifies that Unix() returns the underlying
// instant as seconds since the Unix epoch, including negative values
// for pre-1970 timestamps and exactly zero at the epoch itself.
func TestOffsetDateTime_Unix(t *testing.T) {
	v := NewOffsetDateTime(time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC))
	if got := v.Unix(); got != -14182940 {
		t.Errorf("Expected -14182940, got %d", got)
	}
	epoch := NewOffsetDateTime(time.Unix(0, 0).UTC())
	if got := epoch.Unix(); got != 0 {
		t.Errorf("Expected 0 at Unix epoch, got %d", got)
	}
}

func TestOffsetDateTime_UnmarshalZSuffix(t *testing.T) {
	cases := []struct {
		in       string
		wantSec  int
		wantNsec int
	}{
		{`"2024-07-11T04:37:00Z"`, 0, 0},
		{`"2024-07-11T04:37:00.000Z"`, 0, 0},
		{`"2024-07-11T04:37:00.5Z"`, 0, 500_000_000},
		{`"2024-07-11T04:37:00.123456789Z"`, 0, 123456789},
		{`"2024-07-11T04:37:42Z"`, 42, 0},
	}
	for _, c := range cases {
		var v OffsetDateTime
		if err := json.Unmarshal([]byte(c.in), &v); err != nil {
			t.Errorf("input %s: %v", c.in, err)
			continue
		}
		got := v.AsTime()
		if got.Year() != 2024 || got.Month() != 7 || got.Day() != 11 ||
			got.Hour() != 4 || got.Minute() != 37 ||
			got.Second() != c.wantSec || got.Nanosecond() != c.wantNsec {
			t.Errorf("input %s: got %v", c.in, got)
		}
		if got.Location() != time.UTC {
			t.Errorf("input %s: expected UTC, got %v", c.in, got.Location())
		}
	}
}

func TestOffsetDateTime_UTCRoundTrip(t *testing.T) {
	src := NewOffsetDateTime(time.Date(2024, 7, 11, 4, 37, 0, 123_000_000, time.UTC))
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var dst OffsetDateTime
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("unmarshal %s: %v", data, err)
	}
	if !dst.AsTime().Equal(src.AsTime()) {
		t.Errorf("round-trip mismatch: %v != %v", dst.AsTime(), src.AsTime())
	}
}

func TestOffsetDateTime_In(t *testing.T) {
	plus4 := time.FixedZone("UTC+04:00", 4*3600)

	src := NewOffsetDateTime(time.Date(2024, 7, 11, 10, 0, 0, 0, plus4))
	utc := src.In(time.UTC)
	if got := utc.ToString(); got != "2024-07-11T06:00:00Z" {
		t.Errorf("plus4 -> UTC: expected '2024-07-11T06:00:00Z', got %q", got)
	}

	back := utc.In(plus4)
	if got := back.ToString(); got != "2024-07-11T10:00:00+04:00" {
		t.Errorf("UTC -> plus4: expected '2024-07-11T10:00:00+04:00', got %q", got)
	}

	if !src.AsTime().Equal(utc.AsTime()) || !utc.AsTime().Equal(back.AsTime()) {
		t.Error("In() must preserve the absolute instant")
	}
}

func TestOffsetDateTimeBeforeAfter(t *testing.T) {
	a := NewOffsetDateTime(time.Date(1969, 7, 20, 19, 0, 0, 0, time.UTC))
	b := NewOffsetDateTime(time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC))
	if !a.Before(b) {
		t.Error("Expected a.Before(b) to be true")
	}
	if a.After(b) {
		t.Error("Expected a.After(b) to be false")
	}
}
