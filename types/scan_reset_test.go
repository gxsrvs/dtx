package types

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Scan must derive Valid from the underlying sql.Null*.Valid field, not from
// reflect.TypeOf(value) == nil. A wrapper reused across rows has to transition
// cleanly from a previously-valid state to empty when the next row is NULL.

func TestScanTransitionValidToNull(t *testing.T) {
	t.Run("NullString", func(t *testing.T) {
		var v NullString
		if err := v.Scan("hello"); err != nil {
			t.Fatal(err)
		}
		if !v.Valid || v.Val != "hello" {
			t.Fatalf("first scan: %+v", v)
		}
		if err := v.Scan(nil); err != nil {
			t.Fatal(err)
		}
		if v.Valid || v.Val != "" {
			t.Errorf("after Scan(nil): expected empty+invalid, got %+v", v)
		}
	})

	t.Run("NullBool", func(t *testing.T) {
		var v NullBool
		if err := v.Scan(true); err != nil {
			t.Fatal(err)
		}
		if !v.Valid || !v.Val {
			t.Fatalf("first scan: %+v", v)
		}
		if err := v.Scan(nil); err != nil {
			t.Fatal(err)
		}
		if v.Valid || v.Val {
			t.Errorf("after Scan(nil): expected false+invalid, got %+v", v)
		}
	})

	t.Run("NullInt16", func(t *testing.T) {
		var v NullInt16
		if err := v.Scan(int64(42)); err != nil {
			t.Fatal(err)
		}
		if !v.Valid || v.Val != 42 {
			t.Fatalf("first scan: %+v", v)
		}
		if err := v.Scan(nil); err != nil {
			t.Fatal(err)
		}
		if v.Valid || v.Val != 0 {
			t.Errorf("after Scan(nil): expected 0+invalid, got %+v", v)
		}
	})

	t.Run("NullInt32", func(t *testing.T) {
		var v NullInt32
		if err := v.Scan(int64(123456)); err != nil {
			t.Fatal(err)
		}
		if !v.Valid || v.Val != 123456 {
			t.Fatalf("first scan: %+v", v)
		}
		if err := v.Scan(nil); err != nil {
			t.Fatal(err)
		}
		if v.Valid || v.Val != 0 {
			t.Errorf("after Scan(nil): expected 0+invalid, got %+v", v)
		}
	})

	t.Run("NullInt64", func(t *testing.T) {
		var v NullInt64
		if err := v.Scan(int64(9999999999)); err != nil {
			t.Fatal(err)
		}
		if !v.Valid || v.Val != 9999999999 {
			t.Fatalf("first scan: %+v", v)
		}
		if err := v.Scan(nil); err != nil {
			t.Fatal(err)
		}
		if v.Valid || v.Val != 0 {
			t.Errorf("after Scan(nil): expected 0+invalid, got %+v", v)
		}
	})

	t.Run("NullFloat", func(t *testing.T) {
		var v NullFloat
		if err := v.Scan(float64(3.14)); err != nil {
			t.Fatal(err)
		}
		if !v.Valid || v.Val != 3.14 {
			t.Fatalf("first scan: %+v", v)
		}
		if err := v.Scan(nil); err != nil {
			t.Fatal(err)
		}
		if v.Valid || v.Val != 0 {
			t.Errorf("after Scan(nil): expected 0+invalid, got %+v", v)
		}
	})

	t.Run("NullDate", func(t *testing.T) {
		src := time.Date(1961, 4, 12, 0, 0, 0, 0, time.UTC)
		var v NullDate
		if err := v.Scan(src); err != nil {
			t.Fatal(err)
		}
		if !v.Valid || !v.Val.AsTime().Equal(src) {
			t.Fatalf("first scan: %+v", v)
		}
		if err := v.Scan(nil); err != nil {
			t.Fatal(err)
		}
		if v.Valid || !v.Val.AsTime().IsZero() {
			t.Errorf("after Scan(nil): expected zero+invalid, got %+v", v)
		}
	})

	t.Run("NullDecimal", func(t *testing.T) {
		var v NullDecimal
		if err := v.Scan("12.34"); err != nil {
			t.Fatal(err)
		}
		if !v.Valid || v.Val.String() != "12.34" {
			t.Fatalf("first scan: %+v", v)
		}
		if err := v.Scan(nil); err != nil {
			t.Fatal(err)
		}
		if v.Valid || !v.Val.Equal(decimal.Decimal{}) {
			t.Errorf("after Scan(nil): expected zero+invalid, got %+v", v)
		}
	})

	t.Run("NullUUID", func(t *testing.T) {
		id := uuid.New()
		var v NullUUID
		if err := v.Scan(id.String()); err != nil {
			t.Fatal(err)
		}
		if !v.Valid || v.Val != id {
			t.Fatalf("first scan: %+v", v)
		}
		if err := v.Scan(nil); err != nil {
			t.Fatal(err)
		}
		if v.Valid || v.Val != (uuid.UUID{}) {
			t.Errorf("after Scan(nil): expected zero+invalid, got %+v", v)
		}
	})
}
