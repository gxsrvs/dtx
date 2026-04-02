package types

import (
	"testing"
)

func TestNullInt32FromString(t *testing.T) {
	s := "123456"
	ni := NullInt32FromString(&s)
	if !ni.Valid || ni.Val != 123456 {
		t.Errorf("Expected 123456, got %v", ni)
	}

	niEmpty := NullInt32FromString(nil)
	if niEmpty.Valid {
		t.Error("Expected invalid NullInt32 from nil")
	}

	niInvalid := NullInt32FromString(new("abc"))
	if niInvalid.Valid {
		t.Error("Expected invalid NullInt32 from invalid string")
	}
}

func TestNullInt32FromNullString(t *testing.T) {
	ns := NewNullString("654321")
	ni := NullInt32FromNullString(ns)
	if !ni.Valid || ni.Val != 654321 {
		t.Errorf("Expected 654321, got %v", ni)
	}

	nsEmpty := NewNullStringEmpty()
	niEmpty := NullInt32FromNullString(nsEmpty)
	if niEmpty.Valid {
		t.Error("Expected invalid NullInt32 from empty NullString")
	}

	nsInvalid := NewNullString("xyz")
	niInvalid := NullInt32FromNullString(nsInvalid)
	if niInvalid.Valid {
		t.Error("Expected invalid NullInt32 from invalid NullString")
	}
}
