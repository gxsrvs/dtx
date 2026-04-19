package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDateFromString(t *testing.T) {
	d, err := DateFromString("1961-04-12")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got := d.ToString(); got != "1961-04-12" {
		t.Errorf("Expected 1961-04-12, got %s", got)
	}

	dRu, err := DateFromString("12.04.1961")
	if err != nil {
		t.Fatalf("Unexpected error (RU): %v", err)
	}
	if got := dRu.ToString(); got != "1961-04-12" {
		t.Errorf("Expected 1961-04-12 from RU input, got %s", got)
	}

	if _, err := DateFromString(""); err == nil {
		t.Error("Expected error from empty string, got nil")
	}

	if _, err := DateFromString("not a date"); err == nil {
		t.Error("Expected error from invalid string, got nil")
	}
}

func TestDateAsTime(t *testing.T) {
	src := time.Date(1961, 4, 12, 0, 0, 0, 0, time.UTC)
	d := NewDate(src)
	if !d.AsTime().Equal(src) {
		t.Errorf("Expected %v, got %v", src, d.AsTime())
	}
}

func TestDateMarshalJSON(t *testing.T) {
	d := NewDate(time.Date(1961, 4, 12, 0, 0, 0, 0, time.UTC))
	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if string(data) != `"1961-04-12"` {
		t.Errorf("Expected \"1961-04-12\", got %s", string(data))
	}
}

func TestDateUnmarshalJSON(t *testing.T) {
	var d Date
	if err := json.Unmarshal([]byte(`"1961-04-12"`), &d); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got := d.ToString(); got != "1961-04-12" {
		t.Errorf("Expected 1961-04-12, got %s", got)
	}

	if err := json.Unmarshal([]byte(`null`), &d); err == nil {
		t.Error("Expected error from JSON null, got nil")
	}

	if err := json.Unmarshal([]byte(`""`), &d); err == nil {
		t.Error("Expected error from empty JSON string, got nil")
	}

	if err := json.Unmarshal([]byte(`"not a date"`), &d); err == nil {
		t.Error("Expected error from invalid JSON string, got nil")
	}
}

func TestDateRoundTripJSON(t *testing.T) {
	src := NewDate(time.Date(1961, 4, 12, 0, 0, 0, 0, time.UTC))
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var dst Date
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if src.ToString() != dst.ToString() {
		t.Errorf("Round-trip mismatch: %s != %s", src.ToString(), dst.ToString())
	}
}

func TestDateValue(t *testing.T) {
	src := time.Date(1961, 4, 12, 0, 0, 0, 0, time.UTC)
	d := NewDate(src)
	v, err := d.Value()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	got, ok := v.(time.Time)
	if !ok {
		t.Fatalf("Expected time.Time, got %T", v)
	}
	if !got.Equal(src) {
		t.Errorf("Expected %v, got %v", src, got)
	}
}

func TestDateScan(t *testing.T) {
	src := time.Date(1961, 4, 12, 0, 0, 0, 0, time.UTC)
	var d Date
	if err := d.Scan(src); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !d.AsTime().Equal(src) {
		t.Errorf("Expected %v, got %v", src, d.AsTime())
	}

	var dNull Date
	if err := dNull.Scan(nil); err == nil {
		t.Error("Expected error from NULL scan, got nil")
	}
}
