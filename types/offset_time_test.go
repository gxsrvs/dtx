package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestOffsetTimeFromString(t *testing.T) {
	o, err := OffsetTimeFromString("11:34:51+03:00")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if o.AsTime().Hour() != 11 || o.AsTime().Minute() != 34 || o.AsTime().Second() != 51 {
		t.Errorf("Unexpected time value: %v", o.AsTime())
	}
	_, off := o.AsTime().Zone()
	if off != 3*3600 {
		t.Errorf("Expected offset +03:00, got %d", off)
	}

	if _, err := OffsetTimeFromString(""); err == nil {
		t.Error("Expected error from empty string, got nil")
	}
	if _, err := OffsetTimeFromString("not a time"); err == nil {
		t.Error("Expected error from invalid string, got nil")
	}
}

func TestOffsetTimeRoundTripJSON(t *testing.T) {
	src := NewOffsetTime(time.Date(0, 1, 1, 11, 34, 51, 0, time.FixedZone("UTC+03:00", 3*3600)))
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"11:34:51+03:00"` {
		t.Errorf("Unexpected JSON: %s", string(data))
	}

	var dst OffsetTime
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !dst.AsTime().Equal(src.AsTime()) {
		t.Errorf("Round-trip mismatch: %v != %v", dst.AsTime(), src.AsTime())
	}
}

func TestOffsetTime_StringAndToString(t *testing.T) {
	v := NewOffsetTime(time.Date(0, 1, 1, 20, 17, 40, 0, time.UTC))
	want := "20:17:40Z"
	if got := v.String(); got != want {
		t.Errorf("Expected %q from String(), got %q", want, got)
	}
	if got := v.ToString(); got != want {
		t.Errorf("Expected %q from ToString(), got %q", want, got)
	}
}

func TestOffsetTimeScanBytes(t *testing.T) {
	var v OffsetTime
	if err := v.Scan([]byte("11:34:51+03:00")); err != nil {
		t.Fatalf("Scan([]byte): %v", err)
	}
	_, off := v.AsTime().Zone()
	if off != 3*3600 {
		t.Errorf("Expected offset +03:00 from []byte, got %d", off)
	}

	var bad OffsetTime
	if err := bad.Scan(42); err == nil {
		t.Error("Expected error from unsupported type, got nil")
	}
	if err := bad.Scan("not-a-time"); err == nil {
		t.Error("Expected error from invalid string, got nil")
	}
	if err := bad.Scan([]byte("not-a-time")); err == nil {
		t.Error("Expected error from invalid []byte, got nil")
	}
}

func TestOffsetTimeUnmarshalRejectsNull(t *testing.T) {
	var v OffsetTime
	if err := json.Unmarshal([]byte(`null`), &v); err == nil {
		t.Error("Expected error from JSON null, got nil")
	}
	if err := json.Unmarshal([]byte(`""`), &v); err == nil {
		t.Error("Expected error from empty JSON string, got nil")
	}
}

func TestOffsetTimeValueScan(t *testing.T) {
	src := time.Date(0, 1, 1, 20, 17, 40, 0, time.UTC)
	o := NewOffsetTime(src)
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

	var dst OffsetTime
	if err := dst.Scan(src); err != nil {
		t.Fatalf("Scan(time.Time): %v", err)
	}
	if !dst.AsTime().Equal(src) {
		t.Errorf("Expected %v, got %v", src, dst.AsTime())
	}

	var fromStr OffsetTime
	if err := fromStr.Scan("11:34:51+03:00"); err != nil {
		t.Fatalf("Scan(string): %v", err)
	}
	_, off := fromStr.AsTime().Zone()
	if off != 3*3600 {
		t.Errorf("Expected offset +03:00 from string, got %d", off)
	}

	var nullDst OffsetTime
	if err := nullDst.Scan(nil); err == nil {
		t.Error("Expected error from NULL scan, got nil")
	}
}
