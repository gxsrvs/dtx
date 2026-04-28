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
	if nt.Val.Hour() != 20 || nt.Val.Minute() != 17 || nt.Val.Second() != 40 {
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
	src := NewNullOffsetTime(
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
	if !dst.Val.Equal(src.Val) {
		t.Errorf("Round-trip mismatch: %v != %v", dst.Val, src.Val)
	}
}

func TestNullOffsetTimeMarshalUTC(t *testing.T) {
	v := NewNullOffsetTime(time.Date(0, 1, 1, 20, 17, 40, 0, time.UTC))
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"20:17:40Z"` {
		t.Errorf("Expected Z suffix for UTC, got %s", string(data))
	}
}

func TestNullOffsetTimeToString(t *testing.T) {
	v := NewNullOffsetTime(time.Date(0, 1, 1, 20, 17, 40, 0, time.UTC))
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
	v := NewNullOffsetTime(time.Date(0, 1, 1, 20, 17, 0, 0, time.UTC))
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
	v, err := NewNullOffsetTime(src).Value()
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
	if !v.Valid || !v.Val.Equal(src) {
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
	if v.Val.Hour() != 4 || v.Val.Minute() != 37 || v.Val.Second() != 0 ||
		v.Val.Nanosecond() != 123_000_000 {
		t.Errorf("Unexpected value: %v", v.Val)
	}
	if v.Val.Location() != time.UTC {
		t.Errorf("Expected UTC, got %v", v.Val.Location())
	}
}
