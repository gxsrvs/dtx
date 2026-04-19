package types

import (
	"encoding/json"
	"testing"
)

func TestNullLocalTimeFromString(t *testing.T) {
	s := "02:56:15"
	n := NullLocalTimeFromString(&s)
	if !n.Valid {
		t.Fatal("Expected valid NullLocalTime")
	}
	want := NewLocalTime(2, 56, 15, 0)
	if n.Val != want {
		t.Errorf("Expected %#v, got %#v", want, n.Val)
	}

	if got := NullLocalTimeFromString(nil); got.Valid {
		t.Error("Expected invalid from nil pointer")
	}

	tz := "02:56:15+03:00"
	if got := NullLocalTimeFromString(&tz); got.Valid {
		t.Error("Expected invalid from TZ-bearing string")
	}

	bogus := "not a time"
	if got := NullLocalTimeFromString(&bogus); got.Valid {
		t.Error("Expected invalid from bogus string")
	}
}

func TestNullLocalTimeJSON(t *testing.T) {
	src := NewNullLocalTime(NewLocalTime(2, 56, 15, 0))
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"02:56:15"` {
		t.Errorf("Unexpected JSON: %s", string(data))
	}

	var dst NullLocalTime
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !dst.Valid || dst.Val != src.Val {
		t.Errorf("Round-trip mismatch: %#v != %#v", src, dst)
	}

	// JSON null → invalid
	var nullDst NullLocalTime
	if err := json.Unmarshal([]byte(`null`), &nullDst); err != nil {
		t.Fatalf("unmarshal null: %v", err)
	}
	if nullDst.Valid {
		t.Error("Expected invalid after unmarshal null")
	}

	// Empty marshal
	empty := NewNullLocalTimeEmpty()
	data, err = json.Marshal(empty)
	if err != nil {
		t.Fatalf("marshal empty: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("Expected null, got %s", string(data))
	}
}

func TestNullLocalTimeToString(t *testing.T) {
	v := NewNullLocalTime(NewLocalTime(2, 56, 15, 0))
	if got := v.ToString(); got != "02:56:15" {
		t.Errorf("Expected canonical, got %s", got)
	}
	empty := NewNullLocalTimeEmpty()
	if got := empty.ToString(); got != "" {
		t.Errorf("Expected empty string, got %q", got)
	}
}

func TestNullLocalTimeIsEmpty(t *testing.T) {
	empty := NewNullLocalTimeEmpty()
	if !empty.IsEmpty() {
		t.Error("Expected IsEmpty=true for empty value")
	}
	v := NewNullLocalTime(NewLocalTime(2, 56, 15, 0))
	if v.IsEmpty() {
		t.Error("Expected IsEmpty=false for valid value")
	}
}

func TestNullLocalTimeValueScan(t *testing.T) {
	src := NewNullLocalTime(NewLocalTime(2, 56, 15, 0))
	v, err := src.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	var dst NullLocalTime
	if err := dst.Scan(v); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if !dst.Valid || dst.Val != src.Val {
		t.Errorf("Round-trip mismatch: %#v != %#v", src, dst)
	}

	// Empty Value → nil
	empty := NewNullLocalTimeEmpty()
	v, err = empty.Value()
	if err != nil {
		t.Fatalf("Value (empty): %v", err)
	}
	if v != nil {
		t.Errorf("Expected nil for empty Value, got %v", v)
	}

	// Scan nil → invalid (no error)
	var nullDst NullLocalTime
	if err := nullDst.Scan(nil); err != nil {
		t.Fatalf("Scan(nil): %v", err)
	}
	if nullDst.Valid {
		t.Error("Expected invalid after Scan(nil)")
	}
}
