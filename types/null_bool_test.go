package types

import (
	"encoding/json"
	"testing"
)

func TestNewNullBool(t *testing.T) {
	v := NewNullBool(true)
	if !v.Valid || v.Val != true {
		t.Errorf("Expected (true, valid), got %v", v)
	}

	empty := NewNullBoolEmpty()
	if empty.Valid {
		t.Error("Expected invalid NullBool from NewNullBoolEmpty")
	}
}

func TestNullBoolFromString(t *testing.T) {
	cases := []struct {
		name    string
		input   *string
		wantVal bool
		wantOK  bool
	}{
		{"true", new("true"), true, true},
		{"false", new("false"), false, true},
		{"1", new("1"), true, true},
		{"0", new("0"), false, true},
		{"nil pointer", nil, false, false},
		{"empty string", new(""), false, false},
		{"null word", new("null"), false, false},
		{"NIL upper", new("NIL"), false, false},
		{"garbage", new("yes-please"), false, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := NullBoolFromString(c.input)
			if got.Valid != c.wantOK {
				t.Errorf("Expected Valid=%v, got %v", c.wantOK, got.Valid)
			}
			if got.Valid && got.Val != c.wantVal {
				t.Errorf("Expected Val=%v, got %v", c.wantVal, got.Val)
			}
		})
	}
}

func TestNullBool_IsEmpty(t *testing.T) {
	v := NewNullBool(true)
	if v.IsEmpty() {
		t.Error("Expected IsEmpty=false for valid NullBool")
	}
	empty := NewNullBoolEmpty()
	if !empty.IsEmpty() {
		t.Error("Expected IsEmpty=true for empty NullBool")
	}
}

func TestNullBool_ToString(t *testing.T) {
	v1 := NewNullBool(true)
	if v1.ToString() != "true" {
		t.Error("Expected 'true'")
	}
	v2 := NewNullBool(false)
	if v2.ToString() != "false" {
		t.Error("Expected 'false'")
	}
	empty := NewNullBoolEmpty()
	if empty.ToString() != "" {
		t.Error("Expected empty string for invalid NullBool")
	}
}

func TestNullBool_Value(t *testing.T) {
	v, err := NewNullBool(true).Value()
	if err != nil || v != true {
		t.Errorf("Expected (true, nil), got (%v, %v)", v, err)
	}
	v, err = NewNullBoolEmpty().Value()
	if err != nil || v != nil {
		t.Errorf("Expected (nil, nil), got (%v, %v)", v, err)
	}
}

func TestNullBool_Scan(t *testing.T) {
	cases := []struct {
		name    string
		input   interface{}
		wantOK  bool
		wantVal bool
	}{
		{"bool true", true, true, true},
		{"bool false", false, true, false},
		{"int64 1", int64(1), true, true},
		{"int64 0", int64(0), true, false},
		{"nil", nil, false, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var v NullBool
			if err := v.Scan(c.input); err != nil {
				t.Fatalf("Scan: %v", err)
			}
			if v.Valid != c.wantOK || (v.Valid && v.Val != c.wantVal) {
				t.Errorf("Expected (%v, %v), got %v", c.wantVal, c.wantOK, v)
			}
		})
	}
}

func TestNullBool_MarshalJSON(t *testing.T) {
	data, _ := json.Marshal(NewNullBool(true))
	if string(data) != "true" {
		t.Errorf("Expected true, got %s", data)
	}
	data, _ = json.Marshal(NewNullBoolEmpty())
	if string(data) != "null" {
		t.Errorf("Expected null, got %s", data)
	}
}

func TestNullBool_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		input   string
		wantOK  bool
		wantVal bool
		wantErr bool
	}{
		{"true", true, true, false},
		{"false", true, false, false},
		{"null", false, false, false},
		{`"true"`, true, true, false},
		{"garbage", false, false, true},
	}
	for _, c := range cases {
		var v NullBool
		err := json.Unmarshal([]byte(c.input), &v)
		if (err != nil) != c.wantErr {
			t.Errorf("input %s: wantErr=%v, got err=%v", c.input, c.wantErr, err)
		}
		if !c.wantErr && (v.Valid != c.wantOK || v.Val != c.wantVal) {
			t.Errorf("input %s: expected (val=%v, ok=%v), got %v", c.input, c.wantVal, c.wantOK, v)
		}
	}
}
