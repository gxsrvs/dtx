package types

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

const sampleUuid = "550e8400-e29b-41d4-a716-446655440000"

func TestNullUuidFromString(t *testing.T) {
	s := sampleUuid
	nu := NullUuidFromString(&s)
	if !nu.Valid || nu.Val.String() != s {
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

	for _, s := range []string{"", "null", "NIL"} {
		if got := NullUuidFromString(&s); got.Valid {
			t.Errorf("Expected invalid NullUuid from %q", s)
		}
	}
}

func TestNewNullUuid(t *testing.T) {
	u, _ := uuid.Parse(sampleUuid)
	v := NewNullUuid(u)
	if !v.Valid || v.Val != u {
		t.Errorf("Expected (%v, valid), got %v", u, v)
	}
	if NewNullUuidEmpty().Valid {
		t.Error("Expected invalid NullUuid from NewNullUuidEmpty")
	}
}

func TestNullUuid_IsEmptyToString(t *testing.T) {
	u, _ := uuid.Parse(sampleUuid)
	v := NewNullUuid(u)
	if v.IsEmpty() || v.ToString() != sampleUuid {
		t.Errorf("Expected (false, %s), got (%v, %q)", sampleUuid, v.IsEmpty(), v.ToString())
	}
	empty := NewNullUuidEmpty()
	if !empty.IsEmpty() || empty.ToString() != "" {
		t.Errorf("Expected (true, ''), got (%v, %q)", empty.IsEmpty(), empty.ToString())
	}
}

func TestNullUuid_Scan(t *testing.T) {
	cases := []struct {
		name   string
		input  interface{}
		wantOK bool
	}{
		{"string", sampleUuid, true},
		{"[]byte", []byte(sampleUuid), true},
		{"nil", nil, false},
		{"invalid string", "not-a-uuid", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var v NullUuid
			if err := v.Scan(c.input); err != nil {
				t.Fatalf("Scan: %v", err)
			}
			if v.Valid != c.wantOK {
				t.Errorf("Expected Valid=%v, got %v", c.wantOK, v.Valid)
			}
			if v.Valid && v.Val.String() != sampleUuid {
				t.Errorf("Expected %s, got %s", sampleUuid, v.Val)
			}
		})
	}
}

func TestNullUuid_JSON(t *testing.T) {
	u, _ := uuid.Parse(sampleUuid)
	data, _ := json.Marshal(NewNullUuid(u))
	if string(data) != `"`+sampleUuid+`"` {
		t.Errorf("Expected %q, got %s", sampleUuid, data)
	}
	data, _ = json.Marshal(NewNullUuidEmpty())
	if string(data) != "null" {
		t.Errorf("Expected null, got %s", data)
	}

	cases := []struct {
		input   string
		wantOK  bool
		wantErr bool
	}{
		{`"` + sampleUuid + `"`, true, false},
		{"null", false, false},
		{`"not-a-uuid"`, false, true},
	}
	for _, c := range cases {
		var v NullUuid
		err := json.Unmarshal([]byte(c.input), &v)
		if (err != nil) != c.wantErr {
			t.Errorf("input %s: wantErr=%v, got err=%v", c.input, c.wantErr, err)
		}
		if !c.wantErr && v.Valid != c.wantOK {
			t.Errorf("input %s: wantOK=%v, got Valid=%v", c.input, c.wantOK, v.Valid)
		}
	}
}
