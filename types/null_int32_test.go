package types

import (
	"encoding/json"
	"testing"
)

func TestNewNullInt32(t *testing.T) {
	if v := NewNullInt32(7); !v.Valid || v.Val != 7 {
		t.Errorf("Expected (7, valid), got %v", v)
	}
	if NewNullInt32Empty().Valid {
		t.Error("Expected invalid NullInt32 from NewNullInt32Empty")
	}
}

func TestNullInt32_IsEmptyToString(t *testing.T) {
	v := NewNullInt32(42)
	if v.IsEmpty() || v.ToString() != "42" {
		t.Errorf("Expected (false, '42'), got (%v, %q)", v.IsEmpty(), v.ToString())
	}
	empty := NewNullInt32Empty()
	if !empty.IsEmpty() || empty.ToString() != "" {
		t.Errorf("Expected (true, ''), got (%v, %q)", empty.IsEmpty(), empty.ToString())
	}
}

func TestNullInt32_Value(t *testing.T) {
	v, err := NewNullInt32(42).Value()
	if err != nil || v != int64(42) {
		t.Errorf("Expected (int64(42), nil), got (%v, %v)", v, err)
	}
	v, err = NewNullInt32Empty().Value()
	if err != nil || v != nil {
		t.Errorf("Expected (nil, nil), got (%v, %v)", v, err)
	}
}

func TestNullInt32_Scan(t *testing.T) {
	cases := []struct {
		name    string
		input   interface{}
		wantOK  bool
		wantVal int32
	}{
		{"int64", int64(42), true, 42},
		{"string", "100", true, 100},
		{"nil", nil, false, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var v NullInt32
			if err := v.Scan(c.input); err != nil {
				t.Fatalf("Scan: %v", err)
			}
			if v.Valid != c.wantOK || (v.Valid && v.Val != c.wantVal) {
				t.Errorf("Expected (%v, %v), got %v", c.wantVal, c.wantOK, v)
			}
		})
	}
}

func TestNullInt32_JSON(t *testing.T) {
	data, _ := json.Marshal(NewNullInt32(42))
	if string(data) != "42" {
		t.Errorf("Expected 42, got %s", data)
	}
	data, _ = json.Marshal(NewNullInt32Empty())
	if string(data) != "null" {
		t.Errorf("Expected null, got %s", data)
	}

	cases := []struct {
		input   string
		wantOK  bool
		wantVal int32
		wantErr bool
	}{
		{"42", true, 42, false},
		{`"42"`, true, 42, false},
		{"null", false, 0, false},
		{"abc", false, 0, true},
	}
	for _, c := range cases {
		var v NullInt32
		err := json.Unmarshal([]byte(c.input), &v)
		if (err != nil) != c.wantErr {
			t.Errorf("input %s: wantErr=%v, got err=%v", c.input, c.wantErr, err)
		}
		if !c.wantErr && (v.Valid != c.wantOK || v.Val != c.wantVal) {
			t.Errorf("input %s: expected (val=%v, ok=%v), got %v", c.input, c.wantVal, c.wantOK, v)
		}
	}
}

func TestNullInt32FromString(t *testing.T) {
	s := "123456"
	ni := NullInt32FromString(&s)
	if !ni.Valid || ni.Val != 123456 {
		t.Errorf("Expected 123456, got %v", ni)
	}

	niEmpty := NullInt32FromString(nil)
	if niEmpty.Valid {
		t.Error("Expected invalid NullInt32 from nil")
	}

	niInvalid := NullInt32FromString(new("abc"))
	if niInvalid.Valid {
		t.Error("Expected invalid NullInt32 from invalid string")
	}
}

func TestNullInt32FromNullString(t *testing.T) {
	ns := NewNullString("654321")
	ni := NullInt32FromNullString(ns)
	if !ni.Valid || ni.Val != 654321 {
		t.Errorf("Expected 654321, got %v", ni)
	}

	nsEmpty := NewNullStringEmpty()
	niEmpty := NullInt32FromNullString(nsEmpty)
	if niEmpty.Valid {
		t.Error("Expected invalid NullInt32 from empty NullString")
	}

	nsInvalid := NewNullString("xyz")
	niInvalid := NullInt32FromNullString(nsInvalid)
	if niInvalid.Valid {
		t.Error("Expected invalid NullInt32 from invalid NullString")
	}
}
