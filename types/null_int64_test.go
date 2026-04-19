package types

import (
	"encoding/json"
	"testing"
)

func TestNewNullInt64(t *testing.T) {
	if v := NewNullInt64(7); !v.Valid || v.Val != 7 {
		t.Errorf("Expected (7, valid), got %v", v)
	}
	if NewNullInt64Empty().Valid {
		t.Error("Expected invalid NullInt64 from NewNullInt64Empty")
	}
}

func TestNullInt64_IsEmptyToString(t *testing.T) {
	v := NewNullInt64(42)
	if v.IsEmpty() || v.ToString() != "42" {
		t.Errorf("Expected (false, '42'), got (%v, %q)", v.IsEmpty(), v.ToString())
	}
	empty := NewNullInt64Empty()
	if !empty.IsEmpty() || empty.ToString() != "" {
		t.Errorf("Expected (true, ''), got (%v, %q)", empty.IsEmpty(), empty.ToString())
	}
}

func TestNullInt64_Value(t *testing.T) {
	v, err := NewNullInt64(42).Value()
	if err != nil || v != int64(42) {
		t.Errorf("Expected (int64(42), nil), got (%v, %v)", v, err)
	}
	v, err = NewNullInt64Empty().Value()
	if err != nil || v != nil {
		t.Errorf("Expected (nil, nil), got (%v, %v)", v, err)
	}
}

func TestNullInt64_Scan(t *testing.T) {
	cases := []struct {
		name    string
		input   interface{}
		wantOK  bool
		wantVal int64
	}{
		{"int64", int64(42), true, 42},
		{"string", "100", true, 100},
		{"nil", nil, false, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var v NullInt64
			if err := v.Scan(c.input); err != nil {
				t.Fatalf("Scan: %v", err)
			}
			if v.Valid != c.wantOK || (v.Valid && v.Val != c.wantVal) {
				t.Errorf("Expected (%v, %v), got %v", c.wantVal, c.wantOK, v)
			}
		})
	}
}

func TestNullInt64_JSON(t *testing.T) {
	data, _ := json.Marshal(NewNullInt64(42))
	if string(data) != "42" {
		t.Errorf("Expected 42, got %s", data)
	}
	data, _ = json.Marshal(NewNullInt64Empty())
	if string(data) != "null" {
		t.Errorf("Expected null, got %s", data)
	}

	cases := []struct {
		input   string
		wantOK  bool
		wantVal int64
		wantErr bool
	}{
		{"42", true, 42, false},
		{`"42"`, true, 42, false},
		{"null", false, 0, false},
		{"abc", false, 0, true},
	}
	for _, c := range cases {
		var v NullInt64
		err := json.Unmarshal([]byte(c.input), &v)
		if (err != nil) != c.wantErr {
			t.Errorf("input %s: wantErr=%v, got err=%v", c.input, c.wantErr, err)
		}
		if !c.wantErr && (v.Valid != c.wantOK || v.Val != c.wantVal) {
			t.Errorf("input %s: expected (val=%v, ok=%v), got %v", c.input, c.wantVal, c.wantOK, v)
		}
	}
}

func TestNullInt64FromString(t *testing.T) {
	ni := NullInt64FromString(new("12345"))
	if !ni.Valid || ni.Val != 12345 {
		t.Errorf("Expected 12345, got %v", ni)
	}

	niEmpty := NullInt64FromString(nil)
	if niEmpty.Valid {
		t.Error("Expected invalid NullInt64 from nil")
	}

	niInvalid := NullInt64FromString(new("abc"))
	if niInvalid.Valid {
		t.Error("Expected invalid NullInt64 from invalid string")
	}
}

func TestNullInt64FromNullString(t *testing.T) {
	ns := NewNullString("54321")
	ni := NullInt64FromNullString(ns)
	if !ni.Valid || ni.Val != 54321 {
		t.Errorf("Expected 54321, got %v", ni)
	}

	nsEmpty := NewNullStringEmpty()
	niEmpty := NullInt64FromNullString(nsEmpty)
	if niEmpty.Valid {
		t.Error("Expected invalid NullInt64 from empty NullString")
	}

	nsInvalid := NewNullString("xyz")
	niInvalid := NullInt64FromNullString(nsInvalid)
	if niInvalid.Valid {
		t.Error("Expected invalid NullInt64 from invalid NullString")
	}
}
