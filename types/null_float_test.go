package types

import (
	"encoding/json"
	"math"
	"testing"
)

func TestNewNullFloat(t *testing.T) {
	v := NewNullFloat(3.14)
	if !v.Valid || v.Val != 3.14 {
		t.Errorf("Expected (3.14, valid), got %v", v)
	}
	if NewNullFloatEmpty().Valid {
		t.Error("Expected invalid NullFloat from NewNullFloatEmpty")
	}
}

func TestNullFloatFromString(t *testing.T) {
	cases := []struct {
		name    string
		input   *string
		wantOK  bool
		wantVal float64
	}{
		{"valid", new("3.14"), true, 3.14},
		{"negative", new("-2.5"), true, -2.5},
		{"int as float", new("42"), true, 42.0},
		{"nil pointer", nil, false, 0},
		{"empty string", new(""), false, 0},
		{"null word", new("null"), false, 0},
		{"NIL upper", new("NIL"), false, 0},
		{"garbage", new("not-a-number"), false, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := NullFloatFromString(c.input)
			if got.Valid != c.wantOK {
				t.Errorf("Expected Valid=%v, got %v", c.wantOK, got.Valid)
			}
			if got.Valid && got.Val != c.wantVal {
				t.Errorf("Expected Val=%v, got %v", c.wantVal, got.Val)
			}
		})
	}
}

func TestNullFloat_IsEmpty(t *testing.T) {
	v := NewNullFloat(0)
	if v.IsEmpty() {
		t.Error("Expected IsEmpty=false for valid NullFloat (even with 0)")
	}
	empty := NewNullFloatEmpty()
	if !empty.IsEmpty() {
		t.Error("Expected IsEmpty=true for empty NullFloat")
	}
}

func TestNullFloat_ToString(t *testing.T) {
	v := NewNullFloat(1.5)
	if got := v.ToString(); got != "1.500000" {
		t.Errorf("Expected '1.500000', got %q", got)
	}
	// Documents the M4 inconsistency: empty NullFloat returns "null", not "".
	empty := NewNullFloatEmpty()
	if got := empty.ToString(); got != "null" {
		t.Errorf("Expected 'null' (M4 pending), got %q", got)
	}
}

func TestNullFloat_Value(t *testing.T) {
	v, err := NewNullFloat(1.25).Value()
	if err != nil || v != 1.25 {
		t.Errorf("Expected (1.25, nil), got (%v, %v)", v, err)
	}
	v, err = NewNullFloatEmpty().Value()
	if err != nil || v != nil {
		t.Errorf("Expected (nil, nil), got (%v, %v)", v, err)
	}
}

func TestNullFloat_Scan(t *testing.T) {
	cases := []struct {
		name    string
		input   interface{}
		wantOK  bool
		wantVal float64
	}{
		{"float64", 2.5, true, 2.5},
		{"int64", int64(7), true, 7.0},
		{"string", "1.5", true, 1.5},
		{"nil", nil, false, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var v NullFloat
			if err := v.Scan(c.input); err != nil {
				t.Fatalf("Scan: %v", err)
			}
			if v.Valid != c.wantOK || (v.Valid && v.Val != c.wantVal) {
				t.Errorf("Expected (%v, %v), got %v", c.wantVal, c.wantOK, v)
			}
		})
	}
}

func TestNullFloat_MarshalJSON(t *testing.T) {
	data, _ := json.Marshal(NewNullFloat(1.5))
	if string(data) != "1.5" {
		t.Errorf("Expected 1.5, got %s", data)
	}
	data, _ = json.Marshal(NewNullFloatEmpty())
	if string(data) != "null" {
		t.Errorf("Expected null, got %s", data)
	}
}

func TestNullFloat_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		input   string
		wantOK  bool
		wantVal float64
		wantErr bool
	}{
		{"3.14", true, 3.14, false},
		{`"3.14"`, true, 3.14, false},
		{"null", false, 0, false},
		{"garbage", false, 0, true},
	}
	for _, c := range cases {
		var v NullFloat
		err := json.Unmarshal([]byte(c.input), &v)
		if (err != nil) != c.wantErr {
			t.Errorf("input %s: wantErr=%v, got err=%v", c.input, c.wantErr, err)
		}
		if !c.wantErr && (v.Valid != c.wantOK || v.Val != c.wantVal) {
			t.Errorf("input %s: expected (val=%v, ok=%v), got %v", c.input, c.wantVal, c.wantOK, v)
		}
	}
}

func TestNullFloat_Boundary(t *testing.T) {
	v := NewNullFloat(math.MaxFloat64)
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal MaxFloat64: %v", err)
	}
	var dst NullFloat
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("unmarshal MaxFloat64: %v", err)
	}
	if dst.Val != math.MaxFloat64 {
		t.Errorf("Round-trip lost precision: %v != %v", dst.Val, math.MaxFloat64)
	}
}
