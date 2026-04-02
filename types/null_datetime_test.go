package types

import (
	"testing"
)

func TestNullDateTimeFromString(t *testing.T) {
	ndt := NullDateTimeFromString(new("2023-05-20T15:30:45"))
	if !ndt.Valid {
		t.Error("Expected valid NullDateTime from valid string")
	}
	if ndt.Val.Year() != 2023 || ndt.Val.Month() != 5 || ndt.Val.Day() != 20 ||
		ndt.Val.Hour() != 15 || ndt.Val.Minute() != 30 || ndt.Val.Second() != 45 {
		t.Errorf("Unexpected datetime value: %v", ndt.Val)
	}

	ndtEmpty := NullDateTimeFromString(nil)
	if ndtEmpty.Valid {
		t.Error("Expected invalid NullDateTime from nil pointer")
	}

	ndtInvalid := NullDateTimeFromString(new("invalid"))
	if ndtInvalid.Valid {
		t.Error("Expected invalid NullDateTime from invalid string")
	}
}
