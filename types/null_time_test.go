package types

import (
	"testing"
)

func TestNullTimeFromString(t *testing.T) {
	nt := NullTimeFromString(new("15:30:45"))
	if !nt.Valid {
		t.Error("Expected valid NullTime from valid string")
	}
	if nt.Val.Hour() != 15 || nt.Val.Minute() != 30 || nt.Val.Second() != 45 {
		t.Errorf("Unexpected time value: %v", nt.Val)
	}

	ntEmpty := NullTimeFromString(nil)
	if ntEmpty.Valid {
		t.Error("Expected invalid NullTime from nil pointer")
	}

	ntNull := NullTimeFromString(new("null"))
	if ntNull.Valid {
		t.Error("Expected invalid NullTime from 'null' string")
	}
}
