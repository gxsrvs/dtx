package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewNullDate(t *testing.T) {
	d := time.Date(1961, 4, 12, 0, 0, 0, 0, time.UTC)
	v := NewNullDate(d)
	if !v.Valid || !v.Val.Equal(d) {
		t.Errorf("Expected (%v, valid), got %v", d, v)
	}
	if NewNullDateEmpty().Valid {
		t.Error("Expected invalid NullDate from NewNullDateEmpty")
	}
}

func TestParseDateFromString(t *testing.T) {
	if _, err := ParseDateFromString("1961-04-12"); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if _, err := ParseDateFromString("12.04.1961"); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if _, err := ParseDateFromString("not-a-date"); err == nil {
		t.Error("Expected error from invalid date")
	}
}

func TestNullDate_IsEmpty(t *testing.T) {
	v := NewNullDate(time.Now())
	if v.IsEmpty() {
		t.Error("Expected IsEmpty=false for valid NullDate")
	}
	empty := NewNullDateEmpty()
	if !empty.IsEmpty() {
		t.Error("Expected IsEmpty=true for empty NullDate")
	}
}

func TestNullDate_ToString(t *testing.T) {
	d := time.Date(1961, 4, 12, 0, 0, 0, 0, time.UTC)
	v := NewNullDate(d)
	if got := v.ToString(); got != "1961-04-12" {
		t.Errorf("Expected '1961-04-12', got %q", got)
	}
	empty := NewNullDateEmpty()
	if got := empty.ToString(); got != "" {
		t.Errorf("Expected empty string, got %q", got)
	}
}

func TestNullDate_Value(t *testing.T) {
	d := time.Date(1961, 4, 12, 0, 0, 0, 0, time.UTC)
	v, err := NewNullDate(d).Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	if tv, ok := v.(time.Time); !ok || !tv.Equal(d) {
		t.Errorf("Expected %v, got %v (ok=%v)", d, v, ok)
	}
	v, err = NewNullDateEmpty().Value()
	if err != nil || v != nil {
		t.Errorf("Expected (nil, nil), got (%v, %v)", v, err)
	}
}

func TestNullDate_Scan(t *testing.T) {
	d := time.Date(1961, 4, 12, 0, 0, 0, 0, time.UTC)

	var v NullDate
	if err := v.Scan(d); err != nil {
		t.Fatalf("Scan(time.Time): %v", err)
	}
	if !v.Valid || !v.Val.Equal(d) {
		t.Errorf("Round-trip mismatch: %v vs %v", v, d)
	}

	var n NullDate
	if err := n.Scan(nil); err != nil {
		t.Fatalf("Scan(nil): %v", err)
	}
	if n.Valid {
		t.Error("Expected invalid after Scan(nil)")
	}
}

func TestNullDate_JSON(t *testing.T) {
	d := time.Date(1961, 4, 12, 0, 0, 0, 0, time.UTC)
	data, _ := json.Marshal(NewNullDate(d))
	if string(data) != `"1961-04-12"` {
		t.Errorf("Expected \"1961-04-12\", got %s", data)
	}
	data, _ = json.Marshal(NewNullDateEmpty())
	if string(data) != "null" {
		t.Errorf("Expected null, got %s", data)
	}

	cases := []struct {
		input   string
		wantOK  bool
		wantErr bool
	}{
		{`"1961-04-12"`, true, false},
		{`"12.04.1961"`, true, false},
		{"null", false, false},
		{`"garbage"`, false, true},
	}
	for _, c := range cases {
		var v NullDate
		err := json.Unmarshal([]byte(c.input), &v)
		if (err != nil) != c.wantErr {
			t.Errorf("input %s: wantErr=%v, got err=%v", c.input, c.wantErr, err)
		}
		if !c.wantErr && v.Valid != c.wantOK {
			t.Errorf("input %s: wantOK=%v, got Valid=%v", c.input, c.wantOK, v.Valid)
		}
	}
}

func TestNullDateFromString(t *testing.T) {
	s := "1961-04-12"
	nd := NullDateFromString(&s)
	if !nd.Valid || nd.Val.Format(time.DateOnly) != s {
		t.Errorf("Expected %s, got %v", s, nd)
	}

	sRu := "12.04.1961"
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
	d := time.Date(1961, 4, 12, 0, 0, 0, 0, time.UTC)
	if s := DateToString(d); s != "1961-04-12" {
		t.Errorf("Expected 1961-04-12, got %s", s)
	}
}
