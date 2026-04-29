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
	v := ndt.Val.AsTime()
	if v.Year() != 1961 || v.Month() != 4 || v.Day() != 12 ||
		v.Hour() != 9 || v.Minute() != 7 || v.Second() != 0 {
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
	src := NullOffsetDateTimeFromTime(time.Date(1961, 4, 12, 9, 7, 0, 0, time.FixedZone("UTC+03:00", 3*3600)))
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
	if !dst.Valid || !dst.Val.AsTime().Equal(src.Val.AsTime()) {
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
	v := NullOffsetDateTimeFromTime(time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC))
	if got := v.ToString(); got != "1969-07-20T20:17:40Z" {
		t.Errorf("Expected canonical RFC3339, got %q", got)
	}
	empty := NewNullOffsetDateTimeEmpty()
	if got := empty.ToString(); got != "" {
		t.Errorf("Expected empty string, got %q", got)
	}
}

func TestNullOffsetDateTime_IsEmpty(t *testing.T) {
	v := NullOffsetDateTimeFromTime(time.Now())
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
	v, err := NullOffsetDateTimeFromTime(src).Value()
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

// TestNullOffsetDateTime_UnixOk verifies that UnixOk distinguishes a
// valid epoch-time (e.g. exactly 0 for the Unix epoch) from a NULL
// wrapper (which yields (0, false)).
func TestNullOffsetDateTime_UnixOk(t *testing.T) {
	v := NullOffsetDateTimeFromTime(time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC))
	if got, ok := v.UnixOk(); !ok || got != -14182940 {
		t.Errorf("Expected (-14182940, true), got (%d, %v)", got, ok)
	}
	if got, ok := NewNullOffsetDateTimeEmpty().UnixOk(); ok || got != 0 {
		t.Errorf("Expected (0, false) for NULL, got (%d, %v)", got, ok)
	}

	// Critical regression: epoch=0 must not be confused with NULL=0.
	epoch := NullOffsetDateTimeFromTime(time.Unix(0, 0).UTC())
	if got, ok := epoch.UnixOk(); !ok || got != 0 {
		t.Errorf("Expected (0, true) at Unix epoch, got (%d, %v)", got, ok)
	}
}

func TestNullOffsetDateTime_BeforeAfter(t *testing.T) {
	a := NullOffsetDateTimeFromTime(time.Date(1969, 7, 20, 19, 0, 0, 0, time.UTC))
	b := NullOffsetDateTimeFromTime(time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC))
	if !a.Before(b) {
		t.Error("Expected a.Before(b) to be true")
	}
	if a.After(b) {
		t.Error("Expected a.After(b) to be false")
	}

	empty := NewNullOffsetDateTimeEmpty()
	// Sortable NULL model: NULL is strictly less than any valid value.
	if !empty.Before(b) {
		t.Error("Expected NULL.Before(valid) to be true (sortable)")
	}
	if empty.After(b) {
		t.Error("Expected NULL.After(valid) to be false (sortable)")
	}
	if !b.After(empty) {
		t.Error("Expected valid.After(NULL) to be true (sortable)")
	}
	if b.Before(empty) {
		t.Error("Expected valid.Before(NULL) to be false (sortable)")
	}
}

func TestNullOffsetDateTimeScan(t *testing.T) {
	src := time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC)
	var v NullOffsetDateTime
	if err := v.Scan(src); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if !v.Valid || !v.Val.AsTime().Equal(src) {
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

func TestNullOffsetDateTime_UnmarshalZSuffix(t *testing.T) {
	var v NullOffsetDateTime
	if err := json.Unmarshal([]byte(`"2024-07-11T04:37:00.123Z"`), &v); err != nil {
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

func TestNullOffsetDateTime_In(t *testing.T) {
	plus4 := time.FixedZone("UTC+04:00", 4*3600)

	v := NullOffsetDateTimeFromTime(time.Date(2024, 7, 11, 10, 0, 0, 0, plus4))
	if got := v.In(time.UTC).ToString(); got != "2024-07-11T06:00:00Z" {
		t.Errorf("valid plus4 -> UTC: got %q", got)
	}

	empty := NewNullOffsetDateTimeEmpty()
	got := empty.In(time.UTC)
	if got.Valid {
		t.Error("Expected NULL to propagate through In()")
	}
}

func TestNewNullOffsetDateTime(t *testing.T) {
	odt := NewOffsetDateTime(time.Date(2024, 7, 11, 10, 0, 0, 0, time.UTC))
	v := NewNullOffsetDateTime(odt)
	if !v.Valid || !v.Val.Equal(odt) {
		t.Errorf("Expected (OffsetDateTime, Valid), got %+v", v)
	}
}
