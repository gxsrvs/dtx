package types

import (
	"testing"
)

func TestNullInt64FromString(t *testing.T) {
	ni := NullInt64FromString(new("12345"))
	if !ni.Valid || ni.Val != 12345 {
		t.Errorf("Expected 12345, got %v", ni)
	}

	niEmpty := NullInt64FromString(nil)
	if niEmpty.Valid {
		t.Error("Expected invalid NullInt64 from nil")
	}

	niInvalid := NullInt64FromString(new("abc"))
	if niInvalid.Valid {
		t.Error("Expected invalid NullInt64 from invalid string")
	}
}

func TestNullInt64FromNullString(t *testing.T) {
	ns := NewNullString("54321")
	ni := NullInt64FromNullString(ns)
	if !ni.Valid || ni.Val != 54321 {
		t.Errorf("Expected 54321, got %v", ni)
	}

	nsEmpty := NewNullStringEmpty()
	niEmpty := NullInt64FromNullString(nsEmpty)
	if niEmpty.Valid {
		t.Error("Expected invalid NullInt64 from empty NullString")
	}

	nsInvalid := NewNullString("xyz")
	niInvalid := NullInt64FromNullString(nsInvalid)
	if niInvalid.Valid {
		t.Error("Expected invalid NullInt64 from invalid NullString")
	}
}
