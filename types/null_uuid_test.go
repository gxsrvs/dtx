package types

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

const sampleUUID = "550e8400-e29b-41d4-a716-446655440000"

func TestNullUUIDFromString(t *testing.T) {
	s := sampleUUID
	nu := NullUUIDFromString(&s)
	if !nu.Valid || nu.Val.String() != s {
		t.Errorf("Expected %s, got %v", s, nu)
	}

	nuEmpty := NullUUIDFromString(nil)
	if nuEmpty.Valid {
		t.Error("Expected invalid NullUUID from nil pointer")
	}

	nuInvalid := NullUUIDFromString(new("invalid-uuid"))
	if nuInvalid.Valid {
		t.Error("Expected invalid NullUUID from invalid string")
	}

	for _, s := range []string{"", "null", "NIL"} {
		if got := NullUUIDFromString(&s); got.Valid {
			t.Errorf("Expected invalid NullUUID from %q", s)
		}
	}
}

func TestNewNullUUID(t *testing.T) {
	u, _ := uuid.Parse(sampleUUID)
	v := NewNullUUID(u)
	if !v.Valid || v.Val != u {
		t.Errorf("Expected (%v, valid), got %v", u, v)
	}
	if NewNullUUIDEmpty().Valid {
		t.Error("Expected invalid NullUUID from NewNullUUIDEmpty")
	}
}

func TestNullUUID_IsEmptyToString(t *testing.T) {
	u, _ := uuid.Parse(sampleUUID)
	v := NewNullUUID(u)
	if v.IsEmpty() || v.ToString() != sampleUUID {
		t.Errorf("Expected (false, %s), got (%v, %q)", sampleUUID, v.IsEmpty(), v.ToString())
	}
	empty := NewNullUUIDEmpty()
	if !empty.IsEmpty() || empty.ToString() != "" {
		t.Errorf("Expected (true, ''), got (%v, %q)", empty.IsEmpty(), empty.ToString())
	}
}

func TestNullUUID_Scan(t *testing.T) {
	cases := []struct {
		name   string
		input  interface{}
		wantOK bool
	}{
		{"string", sampleUUID, true},
		{"[]byte", []byte(sampleUUID), true},
		{"nil", nil, false},
		{"invalid string", "not-a-uuid", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var v NullUUID
			if err := v.Scan(c.input); err != nil {
				t.Fatalf("Scan: %v", err)
			}
			if v.Valid != c.wantOK {
				t.Errorf("Expected Valid=%v, got %v", c.wantOK, v.Valid)
			}
			if v.Valid && v.Val.String() != sampleUUID {
				t.Errorf("Expected %s, got %s", sampleUUID, v.Val)
			}
		})
	}
}

// TestNullUUID_Value verifies the driver.Valuer contract: a valid
// wrapper delegates to uuid.UUID.Value and yields the canonical
// 8-4-4-4-12 hex string, while an empty (NULL) wrapper returns
// (nil, nil).
func TestNullUUID_Value(t *testing.T) {
	u, _ := uuid.Parse(sampleUUID)
	v, err := NewNullUUID(u).Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	if s, ok := v.(string); !ok || s != sampleUUID {
		t.Errorf("Expected %q, got %v (type %T)", sampleUUID, v, v)
	}
	v, err = NewNullUUIDEmpty().Value()
	if err != nil || v != nil {
		t.Errorf("Expected (nil, nil), got (%v, %v)", v, err)
	}
}

func TestNullUUID_JSON(t *testing.T) {
	u, _ := uuid.Parse(sampleUUID)
	data, _ := json.Marshal(NewNullUUID(u))
	if string(data) != `"`+sampleUUID+`"` {
		t.Errorf("Expected %q, got %s", sampleUUID, data)
	}
	data, _ = json.Marshal(NewNullUUIDEmpty())
	if string(data) != "null" {
		t.Errorf("Expected null, got %s", data)
	}

	cases := []struct {
		input   string
		wantOK  bool
		wantErr bool
	}{
		{`"` + sampleUUID + `"`, true, false},
		{"null", false, false},
		{`"not-a-uuid"`, false, true},
	}
	for _, c := range cases {
		var v NullUUID
		err := json.Unmarshal([]byte(c.input), &v)
		if (err != nil) != c.wantErr {
			t.Errorf("input %s: wantErr=%v, got err=%v", c.input, c.wantErr, err)
		}
		if !c.wantErr && v.Valid != c.wantOK {
			t.Errorf("input %s: wantOK=%v, got Valid=%v", c.input, c.wantOK, v.Valid)
		}
	}
}
