package types

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewNullDecimal(t *testing.T) {
	d := decimal.NewFromFloat(3.14)
	v := NewNullDecimal(d)
	if !v.Valid || !v.Val.Equal(d) {
		t.Errorf("Expected (%v, valid), got %v", d, v)
	}
	if NewNullDecimalEmpty().Valid {
		t.Error("Expected invalid NullDecimal from NewNullDecimalEmpty")
	}
}

func TestNullDecimalFromString(t *testing.T) {
	cases := []struct {
		name   string
		input  *string
		wantOK bool
	}{
		{"valid", new("3.14"), true},
		{"negative", new("-2.5"), true},
		{"nil pointer", nil, false},
		{"empty", new(""), false},
		{"null word", new("null"), false},
		{"NIL upper", new("NIL"), false},
		{"garbage", new("not-a-number"), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := NullDecimalFromString(c.input)
			if got.Valid != c.wantOK {
				t.Errorf("Expected Valid=%v, got %v", c.wantOK, got.Valid)
			}
		})
	}
}

func TestNullDecimal_IsEmpty(t *testing.T) {
	v := NewNullDecimal(decimal.NewFromInt(0))
	if v.IsEmpty() {
		t.Error("Expected IsEmpty=false for valid NullDecimal")
	}
	empty := NewNullDecimalEmpty()
	if !empty.IsEmpty() {
		t.Error("Expected IsEmpty=true for empty NullDecimal")
	}
}

func TestNullDecimal_Scan(t *testing.T) {
	cases := []struct {
		name   string
		input  interface{}
		wantOK bool
	}{
		{"string", "3.14", true},
		{"[]byte", []byte("2.5"), true},
		{"nil", nil, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var v NullDecimal
			if err := v.Scan(c.input); err != nil {
				t.Fatalf("Scan: %v", err)
			}
			if v.Valid != c.wantOK {
				t.Errorf("Expected Valid=%v, got %v", c.wantOK, v.Valid)
			}
		})
	}
}

func TestNullDecimal_JSON(t *testing.T) {
	d := decimal.NewFromFloat(0.1).Add(decimal.NewFromFloat(0.2))
	src := NewNullDecimal(d)
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) == "null" || string(data) == "" {
		t.Errorf("Expected number, got %s", data)
	}

	var dst NullDecimal
	if err := json.Unmarshal(data, &dst); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !dst.Valid || !dst.Val.Equal(d) {
		t.Errorf("Round-trip mismatch: %v vs %v", dst, src)
	}

	data, _ = json.Marshal(NewNullDecimalEmpty())
	if string(data) != "null" {
		t.Errorf("Expected null, got %s", data)
	}

	var v NullDecimal
	if err := json.Unmarshal([]byte("null"), &v); err != nil {
		t.Fatalf("unmarshal null: %v", err)
	}
	if v.Valid {
		t.Error("Expected invalid after unmarshal null")
	}

	if err := json.Unmarshal([]byte(`"garbage"`), &v); err == nil {
		t.Error("Expected error from invalid input")
	}
}

func TestNullDecimalToString(t *testing.T) {
	d := decimal.NewFromFloat(123.45)
	nd := NewNullDecimal(d)
	if s := nd.ToString(); s != "123.45" {
		t.Errorf("Expected 123.45, got %s", s)
	}

	ndEmpty := NewNullDecimalEmpty()
	if s := ndEmpty.ToString(); s != "" {
		t.Errorf("Expected null, got %s", s)
	}
}

func TestMulNullDecimals(t *testing.T) {
	d1 := decimal.NewFromInt(10)
	d2 := decimal.NewFromInt(5)
	nd1 := NewNullDecimal(d1)
	nd2 := NewNullDecimal(d2)

	res := MulNullDecimals(nd1, nd2)
	if !res.Valid || !res.Val.Equal(decimal.NewFromInt(50)) {
		t.Errorf("Expected 50, got %v", res)
	}

	ndEmpty := NewNullDecimalEmpty()
	resEmpty := MulNullDecimals(nd1, ndEmpty)
	if resEmpty.Valid {
		t.Error("Expected invalid result when multiplying by invalid decimal")
	}
}
