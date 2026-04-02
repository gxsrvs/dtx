package types

import (
	"github.com/shopspring/decimal"
)

type NullDecimal decimal.NullDecimal

func NewNullDecimalEmpty() decimal.NullDecimal {
	return decimal.NullDecimal{
		Decimal: decimal.Decimal{},
		Valid:   false,
	}
}

func NewNullDecimal(val decimal.Decimal) decimal.NullDecimal {
	return decimal.NullDecimal{
		Decimal: val,
		Valid:   true,
	}
}

func NullDecimalToString(val decimal.NullDecimal) string {
	if !val.Valid {
		return "null"
	}
	return val.Decimal.String()
}

func MulNullDecimals(val1, val2 decimal.NullDecimal) decimal.NullDecimal {
	if !val1.Valid || !val2.Valid {
		return NewNullDecimalEmpty()
	}
	return NewNullDecimal(val1.Decimal.Mul(val2.Decimal))
}
