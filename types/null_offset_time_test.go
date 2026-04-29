package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNullOffsetTimeFromString(t *testing.T) {
	nt := NullOffsetTimeFromString(new("20:17:40"))
	if !nt.Valid {
		t.Error("Expected valid NullOffsetTime from valid string")
	}
	v := nt.Val.AsTime()
	if v.Hour() != 20 || v.Minute() != 17 || v.Second() != 40 {
		t.Errorf("Unexpected time value: %v", nt.Val)
	}

	ntEmpty := NullOffsetTimeFromString(nil)
	if ntEmpty.Valid {
		t.Error("Expected invalid NullOffsetTime from nil pointer")
	}

	ntNull := NullOffsetTimeFromString(new("null"))
	if ntNull.Valid {
		t.Error("Expected invalid NullOffsetTime from 'null' string")
	}
}

func TestNullOffsetTimeRoundTripWithTZ(t *testing.T) {
	src := NullOffsetTimeFromTime(
		time.Date(0, 1, 1, 20, 17, 40, 0, time.FixedZone("UTC+03:00", 3*3600)),
	)
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"20:17:40+03:00"` {
		t.Errorf("Unexpected JSON: %s", string(data))
	}

	var dst NullOffsetTime
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !dst.Valid {
		t.Fatal("Expected valid after round-trip")
	}
	if !dst.Val.AsTime().Equal(src.Val.AsTime()) {
		t.Errorf("Round-trip mismatch: %v != %v", dst.Val, src.Val)
	}
}

func TestNullOffsetTimeMarshalUTC(t *testing.T) {
	v := NullOffsetTimeFromTime(time.Date(0, 1, 1, 20, 17, 40, 0, time.UTC))
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"20:17:40Z"` {
		t.Errorf("Expected Z suffix for UTC, got %s", string(data))
	}
}

func TestNullOffsetTimeToString(t *testing.T) {
	v := NullOffsetTimeFromTime(time.Date(0, 1, 1, 20, 17, 40, 0, time.UTC))
	if got := v.ToString(); got != "20:17:40Z" {
		t.Errorf("Expected 20:17:40Z, got %q", got)
	}
	empty := NewNullOffsetTimeEmpty()
	if got := empty.ToString(); got != "" {
		t.Errorf("Expected empty string for invalid value, got %q", got)
	}
}

func TestNullOffsetTimeMarshalEmpty(t *testing.T) {
	v := NewNullOffsetTimeEmpty()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("Expected null, got %s", string(data))
	}
}

func TestNullOffsetTimeUnmarshalNull(t *testing.T) {
	var v NullOffsetTime
	if err := json.Unmarshal([]byte("null"), &v); err != nil {
		t.Fatalf("unmarshal null: %v", err)
	}
	if v.Valid {
		t.Error("Expected invalid after unmarshal null")
	}
}

func TestNullOffsetTime_IsEmpty(t *testing.T) {
	v := NullOffsetTimeFromTime(time.Date(0, 1, 1, 20, 17, 0, 0, time.UTC))
	if v.IsEmpty() {
		t.Error("Expected IsEmpty=false for valid NullOffsetTime")
	}
	empty := NewNullOffsetTimeEmpty()
	if !empty.IsEmpty() {
		t.Error("Expected IsEmpty=true for empty NullOffsetTime")
	}
}

func TestNullOffsetTime_Value(t *testing.T) {
	src := time.Date(0, 1, 1, 20, 17, 40, 0, time.UTC)
	v, err := NullOffsetTimeFromTime(src).Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	if tv, ok := v.(time.Time); !ok || !tv.Equal(src) {
		t.Errorf("Expected %v, got %v (ok=%v)", src, v, ok)
	}
	v, err = NewNullOffsetTimeEmpty().Value()
	if err != nil || v != nil {
		t.Errorf("Expected (nil, nil), got (%v, %v)", v, err)
	}
}

func TestNullOffsetTimeFromStringInvalid(t *testing.T) {
	if got := NullOffsetTimeFromString(new("not-a-time")); got.Valid {
		t.Error("Expected invalid NullOffsetTime from garbage")
	}
}

func TestNullOffsetTimeScan(t *testing.T) {
	src := time.Date(0, 1, 1, 20, 17, 40, 0, time.UTC)
	var v NullOffsetTime
	if err := v.Scan(src); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if !v.Valid || !v.Val.AsTime().Equal(src) {
		t.Errorf("Round-trip mismatch: %#v vs %v", v, src)
	}

	var nullV NullOffsetTime
	if err := nullV.Scan(nil); err != nil {
		t.Fatalf("Scan(nil): %v", err)
	}
	if nullV.Valid {
		t.Error("Expected invalid after Scan(nil)")
	}
}

func TestNullOffsetTime_UnmarshalZSuffix(t *testing.T) {
	var v NullOffsetTime
	if err := json.Unmarshal([]byte(`"04:37:00.123Z"`), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !v.Valid {
		t.Fatal("Expected valid")
	}
	got := v.Val.AsTime()
	if got.Hour() != 4 || got.Minute() != 37 || got.Second() != 0 ||
		got.Nanosecond() != 123_000_000 {
		t.Errorf("Unexpected value: %v", v.Val)
	}
	if got.Location() != time.UTC {
		t.Errorf("Expected UTC, got %v", got.Location())
	}
}

func TestNullOffsetTime_In(t *testing.T) {
	plus4 := time.FixedZone("UTC+04:00", 4*3600)

	v := NullOffsetTimeFromTime(time.Date(0, 1, 1, 10, 0, 0, 0, plus4))
	if got := v.In(time.UTC).ToString(); got != "06:00:00Z" {
		t.Errorf("valid plus4 -> UTC: got %q", got)
	}

	empty := NewNullOffsetTimeEmpty()
	got := empty.In(time.UTC)
	if got.Valid {
		t.Error("Expected NULL to propagate through In()")
	}
}

func TestNewNullOffsetTime(t *testing.T) {
	ot := NewOffsetTime(time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC))
	v := NewNullOffsetTime(ot)
	if !v.Valid || !v.Val.Equal(ot) {
		t.Errorf("Expected (OffsetTime, Valid), got %+v", v)
	}
}
