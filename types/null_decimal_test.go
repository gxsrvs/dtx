package types

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestNullDecimalToString(t *testing.T) {
	d := decimal.NewFromFloat(123.45)
	nd := NewNullDecimal(d)
	if s := NullDecimalToString(nd); s != "123.45" {
		t.Errorf("Expected 123.45, got %s", s)
	}

	ndEmpty := NewNullDecimalEmpty()
	if s := NullDecimalToString(ndEmpty); s != "null" {
		t.Errorf("Expected null, got %s", s)
	}
}

func TestMulNullDecimals(t *testing.T) {
	d1 := decimal.NewFromInt(10)
	d2 := decimal.NewFromInt(5)
	nd1 := NewNullDecimal(d1)
	nd2 := NewNullDecimal(d2)

	res := MulNullDecimals(nd1, nd2)
	if !res.Valid || !res.Decimal.Equal(decimal.NewFromInt(50)) {
		t.Errorf("Expected 50, got %v", res)
	}

	ndEmpty := NewNullDecimalEmpty()
	resEmpty := MulNullDecimals(nd1, ndEmpty)
	if resEmpty.Valid {
		t.Error("Expected invalid result when multiplying by invalid decimal")
	}
}
