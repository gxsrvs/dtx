package types

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/shopspring/decimal"
)

type NullDecimal struct {
	Val   decimal.Decimal
	Valid bool // Valid is true if UUID is not NULL
}

func NewNullDecimal(val decimal.Decimal) NullDecimal {
	return NullDecimal{
		Val:   val,
		Valid: true,
	}
}

func NewNullDecimalEmpty() NullDecimal {
	return NullDecimal{
		Val:   decimal.Decimal{},
		Valid: false,
	}
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDecimal) IsEmpty() bool {
	return !thisVal.Valid
}

//goland:noinspection GoUnusedExportedFunction
func (thisVal *NullDecimal) ToString() string {
	if !thisVal.Valid {
		return ""
	}
	return thisVal.Val.String()
}

func NullDecimalFromString(strValue *string) NullDecimal {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullDecimalEmpty()
	}
	result, err := decimal.NewFromString(*strValue)
	if err != nil {
		return NewNullDecimalEmpty()
	}
	return NewNullDecimal(result)
}

func MulNullDecimals(val1, val2 NullDecimal) NullDecimal {
	if !val1.Valid || !val2.Valid {
		return NewNullDecimalEmpty()
	}
	return NewNullDecimal(val1.Val.Mul(val2.Val))
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDecimal) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*thisVal = NewNullDecimalEmpty()
	} else {
		*thisVal = NullDecimalFromString(&s.String)
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (thisVal NullDecimal) MarshalJSON() ([]byte, error) {
	if !thisVal.Valid {
		return nullJson, nil
	}
	return json.Marshal(thisVal.Val)

}

//goland:noinspection GoMixedReceiverTypes
func (thisVal *NullDecimal) UnmarshalJSON(data []byte) error {
	sd := string(data)
	if sd == "null" || sd == "" {
		thisVal.Valid = false
		thisVal.Val = decimal.Decimal{}
		return nil
	}
	err := json.Unmarshal(data, &thisVal.Val)
	if err != nil {
		thisVal.Valid = false
		return err
	}
	thisVal.Valid = true
	return nil
}
