package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNullLocalDateTimeFromString(t *testing.T) {
	s := "1969-07-21T02:56:15"
	n := NullLocalDateTimeFromString(&s)
	if !n.Valid {
		t.Fatal("Expected valid NullLocalDateTime")
	}
	want := NewLocalDateTime(1969, time.July, 21, 2, 56, 15, 0)
	if n.Val != want {
		t.Errorf("Expected %#v, got %#v", want, n.Val)
	}

	if got := NullLocalDateTimeFromString(nil); got.Valid {
		t.Error("Expected invalid from nil pointer")
	}

	tz := "1969-07-21T02:56:15+03:00"
	if got := NullLocalDateTimeFromString(&tz); got.Valid {
		t.Error("Expected invalid from TZ-bearing string")
	}

	bogus := "not a datetime"
	if got := NullLocalDateTimeFromString(&bogus); got.Valid {
		t.Error("Expected invalid from bogus string")
	}
}

func TestNullLocalDateTimeJSON(t *testing.T) {
	src := NewNullLocalDateTime(NewLocalDateTime(1969, time.July, 21, 2, 56, 15, 0))
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"1969-07-21T02:56:15"` {
		t.Errorf("Unexpected JSON: %s", string(data))
	}

	var dst NullLocalDateTime
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !dst.Valid || dst.Val != src.Val {
		t.Errorf("Round-trip mismatch: %#v != %#v", src, dst)
	}

	// JSON null → invalid (not an error for nullable)
	var nullDst NullLocalDateTime
	if err := json.Unmarshal([]byte(`null`), &nullDst); err != nil {
		t.Fatalf("unmarshal null: %v", err)
	}
	if nullDst.Valid {
		t.Error("Expected invalid after unmarshal null")
	}

	// Empty marshal
	empty := NewNullLocalDateTimeEmpty()
	data, err = json.Marshal(empty)
	if err != nil {
		t.Fatalf("marshal empty: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("Expected null, got %s", string(data))
	}
}

func TestNullLocalDateTimeToString(t *testing.T) {
	v := NewNullLocalDateTime(NewLocalDateTime(1969, time.July, 21, 2, 56, 15, 0))
	if got := v.ToString(); got != "1969-07-21T02:56:15" {
		t.Errorf("Expected canonical, got %s", got)
	}
	empty := NewNullLocalDateTimeEmpty()
	if got := empty.ToString(); got != "" {
		t.Errorf("Expected empty string, got %q", got)
	}
}

func TestNullLocalDateTimeIsEmpty(t *testing.T) {
	empty := NewNullLocalDateTimeEmpty()
	if !empty.IsEmpty() {
		t.Error("Expected IsEmpty=true for empty value")
	}
	v := NewNullLocalDateTime(NewLocalDateTime(1969, time.July, 21, 2, 56, 15, 0))
	if v.IsEmpty() {
		t.Error("Expected IsEmpty=false for valid value")
	}
}

func TestNullLocalDateTimeValueScan(t *testing.T) {
	src := NewNullLocalDateTime(NewLocalDateTime(1969, time.July, 21, 2, 56, 15, 0))
	v, err := src.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	var dst NullLocalDateTime
	if err := dst.Scan(v); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if !dst.Valid || dst.Val != src.Val {
		t.Errorf("Round-trip mismatch: %#v != %#v", src, dst)
	}

	// Empty Value → nil
	empty := NewNullLocalDateTimeEmpty()
	v, err = empty.Value()
	if err != nil {
		t.Fatalf("Value (empty): %v", err)
	}
	if v != nil {
		t.Errorf("Expected nil for empty Value, got %v", v)
	}

	// Scan nil → invalid (no error)
	var nullDst NullLocalDateTime
	if err := nullDst.Scan(nil); err != nil {
		t.Fatalf("Scan(nil): %v", err)
	}
	if nullDst.Valid {
		t.Error("Expected invalid after Scan(nil)")
	}
}

func TestNullLocalDateTimeFromTime(t *testing.T) {
	src := time.Date(2024, 7, 11, 10, 30, 45, 0, time.UTC)
	v := NullLocalDateTimeFromTime(src)
	if !v.Valid {
		t.Fatal("Expected valid NullLocalDateTime")
	}
	if v.Val.Year != 2024 || v.Val.Month != time.July || v.Val.Day != 11 ||
		v.Val.Hour != 10 || v.Val.Minute != 30 || v.Val.Second != 45 {
		t.Errorf("Unexpected fields: %+v", v.Val)
	}
}
