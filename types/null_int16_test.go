package types

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"
)

func TestNullInt16FromString(t *testing.T) {
	ni := NullInt16FromString(new("123"))
	if !ni.Valid || ni.Val != 123 {
		t.Errorf("Expected 123, got %v", ni)
	}

	niEmpty := NullInt16FromString(nil)
	if niEmpty.Valid {
		t.Error("Expected invalid NullInt16 from nil")
	}

	niInvalid := NullInt16FromString(new("abc"))
	if niInvalid.Valid {
		t.Error("Expected invalid NullInt16 from invalid string")
	}

	// Overflow → invalid
	overflow := strconv.Itoa(math.MaxInt32)
	if got := NullInt16FromString(&overflow); got.Valid {
		t.Error("Expected invalid NullInt16 from overflow")
	}

	// "null"/"nil" → invalid
	for _, s := range []string{"null", "NIL", ""} {
		if got := NullInt16FromString(&s); got.Valid {
			t.Errorf("Expected invalid NullInt16 from %q", s)
		}
	}
}

func TestNullInt16FromNullString(t *testing.T) {
	if got := NullInt16FromNullString(NewNullString("777")); !got.Valid || got.Val != 777 {
		t.Errorf("Expected 777, got %v", got)
	}
	if got := NullInt16FromNullString(NewNullStringEmpty()); got.Valid {
		t.Error("Expected invalid NullInt16 from empty NullString")
	}
	if got := NullInt16FromNullString(NewNullString("abc")); got.Valid {
		t.Error("Expected invalid NullInt16 from non-numeric NullString")
	}
}

func TestNewNullInt16(t *testing.T) {
	if v := NewNullInt16(7); !v.Valid || v.Val != 7 {
		t.Errorf("Expected (7, valid), got %v", v)
	}
	if NewNullInt16Empty().Valid {
		t.Error("Expected invalid NullInt16 from NewNullInt16Empty")
	}
}

func TestNullInt16_IsEmptyToString(t *testing.T) {
	v := NewNullInt16(42)
	if v.IsEmpty() || v.ToString() != "42" {
		t.Errorf("Expected (false, '42'), got (%v, %q)", v.IsEmpty(), v.ToString())
	}
	empty := NewNullInt16Empty()
	if !empty.IsEmpty() || empty.ToString() != "" {
		t.Errorf("Expected (true, ''), got (%v, %q)", empty.IsEmpty(), empty.ToString())
	}
}

func TestNullInt16_Value(t *testing.T) {
	v, err := NewNullInt16(42).Value()
	if err != nil || v != int64(42) {
		t.Errorf("Expected (int64(42), nil), got (%v, %v)", v, err)
	}
	v, err = NewNullInt16Empty().Value()
	if err != nil || v != nil {
		t.Errorf("Expected (nil, nil), got (%v, %v)", v, err)
	}
}

func TestNullInt16_Scan(t *testing.T) {
	cases := []struct {
		name    string
		input   interface{}
		wantOK  bool
		wantVal int16
	}{
		{"int64", int64(42), true, 42},
		{"int16-fits-string", "100", true, 100},
		{"nil", nil, false, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var v NullInt16
			if err := v.Scan(c.input); err != nil {
				t.Fatalf("Scan: %v", err)
			}
			if v.Valid != c.wantOK || (v.Valid && v.Val != c.wantVal) {
				t.Errorf("Expected (%v, %v), got %v", c.wantVal, c.wantOK, v)
			}
		})
	}
}

func TestNullInt16_JSON(t *testing.T) {
	data, _ := json.Marshal(NewNullInt16(42))
	if string(data) != "42" {
		t.Errorf("Expected 42, got %s", data)
	}
	data, _ = json.Marshal(NewNullInt16Empty())
	if string(data) != "null" {
		t.Errorf("Expected null, got %s", data)
	}

	cases := []struct {
		input   string
		wantOK  bool
		wantVal int16
		wantErr bool
	}{
		{"42", true, 42, false},
		{`"42"`, true, 42, false},
		{"null", false, 0, false},
		{"abc", false, 0, true},
	}
	for _, c := range cases {
		var v NullInt16
		err := json.Unmarshal([]byte(c.input), &v)
		if (err != nil) != c.wantErr {
			t.Errorf("input %s: wantErr=%v, got err=%v", c.input, c.wantErr, err)
		}
		if !c.wantErr && (v.Valid != c.wantOK || v.Val != c.wantVal) {
			t.Errorf("input %s: expected (val=%v, ok=%v), got %v", c.input, c.wantVal, c.wantOK, v)
		}
	}
}
