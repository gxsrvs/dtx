package types

import (
	"testing"
	"time"
)

func TestNullDateFromString(t *testing.T) {
	s := "2023-05-20"
	nd := NullDateFromString(&s)
	if !nd.Valid || nd.Val.Format(time.DateOnly) != s {
		t.Errorf("Expected %s, got %v", s, nd)
	}

	sRu := "20.05.2023"
	ndRu := NullDateFromString(&sRu)
	if !ndRu.Valid || ndRu.Val.Format(RuOnlyDateMask) != sRu {
		t.Errorf("Expected %s (RU), got %v", sRu, ndRu)
	}

	ndEmpty := NullDateFromString(nil)
	if ndEmpty.Valid {
		t.Error("Expected invalid NullDate from nil pointer")
	}

	ndInvalid := NullDateFromString(new("invalid"))
	if ndInvalid.Valid {
		t.Error("Expected invalid NullDate from invalid string")
	}
}

func TestDateToString(t *testing.T) {
	d := time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC)
	if s := DateToString(d); s != "2023-05-20" {
		t.Errorf("Expected 2023-05-20, got %s", s)
	}
}
