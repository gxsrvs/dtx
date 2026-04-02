package types

import (
	"testing"
)

func TestNullUuidFromString(t *testing.T) {
	s := "550e8400-e29b-41d4-a716-446655440000"
	nu := NullUuidFromString(&s)
	if !nu.Valid || nu.UUID.String() != s {
		t.Errorf("Expected %s, got %v", s, nu)
	}

	nuEmpty := NullUuidFromString(nil)
	if nuEmpty.Valid {
		t.Error("Expected invalid NullUuid from nil pointer")
	}

	nuInvalid := NullUuidFromString(new("invalid-uuid"))
	if nuInvalid.Valid {
		t.Error("Expected invalid NullUuid from invalid string")
	}
}
