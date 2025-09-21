package types

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type NullDecimal decimal.NullDecimal

//goland:noinspection GoMixedReceiverTypes,GoUnusedExportedFunction
func NullDecimalToString(val decimal.NullDecimal) string {
	if !val.Valid {
		return "null"
	}
	return fmt.Sprintf("%f", val.Decimal)
}
