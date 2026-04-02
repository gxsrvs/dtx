package types

import (
	"testing"
	"time"
)

func TestMaxDateTime(t *testing.T) {
	t1 := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)

	if res := MaxDateTime(t1, t2); !res.Equal(t2) {
		t.Errorf("Expected %v, got %v", t2, res)
	}
	if res := MaxDateTime(t2, t1); !res.Equal(t2) {
		t.Errorf("Expected %v, got %v", t2, res)
	}
}

func TestAssembleNullDateTimeTZ(t *testing.T) {
	dateVal := time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC)
	timeVal := time.Date(0, 1, 1, 15, 30, 45, 0, time.UTC)

	nd := NewNullDate(dateVal)
	nt := NewNullTime(timeVal)

	tz := "+03:00"
	res, err := AssembleNullDateTimeTZ(&nd, nil, &nt, nil, tz)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedLoc := time.FixedZone("UTC+03:00", 3*3600)
	expected := time.Date(2023, 5, 20, 15, 30, 45, 0, expectedLoc)

	if !res.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, res)
	}

	// Test with defaults
	resDefault, err := AssembleNullDateTimeTZ(new(NewNullDateEmpty()), &dateVal, new(NewNullTimeEmpty()), &timeVal, tz)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !resDefault.Equal(expected) {
		t.Errorf("Expected %v (default), got %v", expected, resDefault)
	}
}
