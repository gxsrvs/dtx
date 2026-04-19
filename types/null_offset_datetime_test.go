package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNullOffsetDateTimeFromString(t *testing.T) {
	ndt := NullOffsetDateTimeFromString(new("1961-04-12T09:07:00"))
	if !ndt.Valid {
		t.Error("Expected valid NullOffsetDateTime from valid string")
	}
	if ndt.Val.Year() != 1961 || ndt.Val.Month() != 4 || ndt.Val.Day() != 12 ||
		ndt.Val.Hour() != 9 || ndt.Val.Minute() != 7 || ndt.Val.Second() != 0 {
		t.Errorf("Unexpected datetime value: %v", ndt.Val)
	}

	if got := NullOffsetDateTimeFromString(nil); got.Valid {
		t.Error("Expected invalid NullOffsetDateTime from nil pointer")
	}

	bogus := "invalid"
	if got := NullOffsetDateTimeFromString(&bogus); got.Valid {
		t.Error("Expected invalid NullOffsetDateTime from invalid string")
	}
}

func TestNullOffsetDateTimeRoundTripJSON(t *testing.T) {
	src := NewNullOffsetDateTime(time.Date(1961, 4, 12, 9, 7, 0, 0, time.FixedZone("UTC+03:00", 3*3600)))
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"1961-04-12T09:07:00+03:00"` {
		t.Errorf("Unexpected JSON: %s", string(data))
	}

	var dst NullOffsetDateTime
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !dst.Valid || !dst.Val.Equal(src.Val) {
		t.Errorf("Round-trip mismatch: %#v != %#v", dst, src)
	}
}

func TestNullOffsetDateTimeMarshalEmpty(t *testing.T) {
	v := NewNullOffsetDateTimeEmpty()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("Expected null, got %s", string(data))
	}
}

func TestNullOffsetDateTimeUnmarshalNull(t *testing.T) {
	var v NullOffsetDateTime
	if err := json.Unmarshal([]byte("null"), &v); err != nil {
		t.Fatalf("unmarshal null: %v", err)
	}
	if v.Valid {
		t.Error("Expected invalid after unmarshal null")
	}
}

func TestNullOffsetDateTimeToString(t *testing.T) {
	v := NewNullOffsetDateTime(time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC))
	if got := v.ToString(); got != "1969-07-20T20:17:40Z" {
		t.Errorf("Expected canonical RFC3339, got %q", got)
	}
	empty := NewNullOffsetDateTimeEmpty()
	if got := empty.ToString(); got != "" {
		t.Errorf("Expected empty string, got %q", got)
	}
}

func TestNullOffsetDateTime_IsEmpty(t *testing.T) {
	v := NewNullOffsetDateTime(time.Now())
	if v.IsEmpty() {
		t.Error("Expected IsEmpty=false for valid NullOffsetDateTime")
	}
	empty := NewNullOffsetDateTimeEmpty()
	if !empty.IsEmpty() {
		t.Error("Expected IsEmpty=true for empty NullOffsetDateTime")
	}
}

func TestNullOffsetDateTime_Value(t *testing.T) {
	src := time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC)
	v, err := NewNullOffsetDateTime(src).Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	if tv, ok := v.(time.Time); !ok || !tv.Equal(src) {
		t.Errorf("Expected %v, got %v (ok=%v)", src, v, ok)
	}
	v, err = NewNullOffsetDateTimeEmpty().Value()
	if err != nil || v != nil {
		t.Errorf("Expected (nil, nil), got (%v, %v)", v, err)
	}
}

func TestNullOffsetDateTime_BeforeAfter(t *testing.T) {
	a := NewNullOffsetDateTime(time.Date(1969, 7, 20, 19, 0, 0, 0, time.UTC))
	b := time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC)
	if !a.Before(b) {
		t.Error("Expected a.Before(b) to be true")
	}
	if a.After(b) {
		t.Error("Expected a.After(b) to be false")
	}

	empty := NewNullOffsetDateTimeEmpty()
	if empty.Before(b) || empty.After(b) {
		t.Error("Expected empty Before/After to return false")
	}
}

func TestNullOffsetDateTimeScan(t *testing.T) {
	src := time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC)
	var v NullOffsetDateTime
	if err := v.Scan(src); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if !v.Valid || !v.Val.Equal(src) {
		t.Errorf("Round-trip mismatch")
	}

	var nullV NullOffsetDateTime
	if err := nullV.Scan(nil); err != nil {
		t.Fatalf("Scan(nil): %v", err)
	}
	if nullV.Valid {
		t.Error("Expected invalid after Scan(nil)")
	}
}
