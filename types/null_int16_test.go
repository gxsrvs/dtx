package types

import (
	"testing"
)

func TestNullInt16FromString(t *testing.T) {
	ni := NullInt16FromString(new("123"))
	if !ni.Valid || ni.Val != 123 {
		t.Errorf("Expected 123, got %v", ni)
	}

	niEmpty := NullInt16FromString(nil)
	if niEmpty.Valid {
		t.Error("Expected invalid NullInt16 from nil")
	}

	niInvalid := NullInt16FromString(new("abc"))
	if niInvalid.Valid {
		t.Error("Expected invalid NullInt16 from invalid string")
	}
}
