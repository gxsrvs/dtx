package types

import (
	"database/sql"
	"testing"
)

func TestNewNullString(t *testing.T) {
	val := "test"
	ns := NewNullString(val)
	if ns.Val != val {
		t.Errorf("expected %s, got %s", val, ns.Val)
	}
	if !ns.Valid {
		t.Error("expected Valid to be true")
	}
}

func TestNewNullStringEmpty(t *testing.T) {
	ns := NewNullStringEmpty()
	if ns.Valid {
		t.Error("expected Valid to be false")
	}
}

func TestNSFromString(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"hello", true},
		{"", false},
	}

	for _, c := range cases {
		ns := NSFromString(c.input)
		if ns.String != c.input {
			t.Errorf("input %s: expected String %s, got %s", c.input, c.input, ns.String)
		}
		if ns.Valid != c.expected {
			t.Errorf("input %s: expected Valid %v, got %v", c.input, c.expected, ns.Valid)
		}
	}
}

func TestGetNullString(t *testing.T) {
	cases := []struct {
		input    sql.NullString
		expected string
	}{
		{sql.NullString{String: "hello", Valid: true}, "hello"},
		{sql.NullString{String: "", Valid: false}, ""},
		{sql.NullString{String: "ignored", Valid: false}, ""},
	}

	for _, c := range cases {
		res := GetNullString(c.input)
		if res != c.expected {
			t.Errorf("expected %s, got %s", c.expected, res)
		}
	}
}

func TestNullString_IsEmpty(t *testing.T) {
	ns1 := NewNullString("val")
	if ns1.IsEmpty() {
		t.Error("expected IsEmpty to be false for valid string")
	}

	ns2 := NewNullStringEmpty()
	if !ns2.IsEmpty() {
		t.Error("expected IsEmpty to be true for empty string")
	}
}

func TestNullString_ToString(t *testing.T) {
	ns1 := NewNullString("val")
	if ns1.ToString() != "val" {
		t.Errorf("expected val, got %s", ns1.ToString())
	}

	ns2 := NewNullStringEmpty()
	if ns2.ToString() != "" {
		t.Errorf("expected empty string, got %s", ns2.ToString())
	}
}

func TestNullString_Value(t *testing.T) {
	ns1 := NewNullString("val")
	val1, err1 := ns1.Value()
	if err1 != nil {
		t.Errorf("unexpected error: %v", err1)
	}
	if val1 != "val" {
		t.Errorf("expected val, got %v", val1)
	}

	ns2 := NewNullStringEmpty()
	val2, err2 := ns2.Value()
	if err2 != nil {
		t.Errorf("unexpected error: %v", err2)
	}
	if val2 != nil {
		t.Errorf("expected nil, got %v", val2)
	}
}

func TestNullString_Scan(t *testing.T) {
	var ns NullString
	err := ns.Scan("hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ns.Val != "hello" || !ns.Valid {
		t.Errorf("expected hello/true, got %s/%v", ns.Val, ns.Valid)
	}

	err = ns.Scan(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ns.Valid {
		t.Error("expected Valid to be false after scanning nil")
	}
}

func TestNullString_JSON(t *testing.T) {
	// Marshal
	ns1 := NewNullString("hello")
	data1, err1 := ns1.MarshalJSON()
	if err1 != nil {
		t.Errorf("marshal error: %v", err1)
	}
	if string(data1) != `"hello"` {
		t.Errorf("expected \"hello\", got %s", string(data1))
	}

	ns2 := NewNullStringEmpty()
	data2, err2 := ns2.MarshalJSON()
	if err2 != nil {
		t.Errorf("marshal error: %v", err2)
	}
	if string(data2) != "null" {
		t.Errorf("expected null, got %s", string(data2))
	}

	// Unmarshal
	var ns3 NullString
	err3 := ns3.UnmarshalJSON([]byte(`"world"`))
	if err3 != nil {
		t.Errorf("unmarshal error: %v", err3)
	}
	if ns3.Val != "world" || !ns3.Valid {
		t.Errorf("expected world/true, got %s/%v", ns3.Val, ns3.Valid)
	}

	var ns4 NullString
	err4 := ns4.UnmarshalJSON([]byte(`null`))
	if err4 != nil {
		t.Errorf("unmarshal error: %v", err4)
	}
	if ns4.Valid {
		t.Error("expected Valid to be false after unmarshaling null")
	}
}

func TestNullString_ReflectTypeOfNil(t *testing.T) {
	// Специальный тест для проверки логики в Scan:
	// if reflect.TypeOf(value) == nil {
	var ns NullString
	err := ns.Scan(nil)
	if err != nil {
		t.Fatal(err)
	}
	if ns.Valid {
		t.Error("Should be invalid for nil")
	}
}
