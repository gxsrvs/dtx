package types

import (
	"testing"
	"time"
)

// TestSortableComparison_NullDate verifies the full sortable NULL
// semantics matrix for NullDate. The other Null* types delegate to
// the same underlying logic; per-type tests below are smoke checks
// that NULL is strictly less than any valid value.
func TestSortableComparison_NullDate(t *testing.T) {
	a := NullDateFromTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	b := NullDateFromTime(time.Date(2024, 7, 11, 0, 0, 0, 0, time.UTC))
	z := NewNullDateEmpty()
	z2 := NewNullDateEmpty()

	// valid vs valid
	if !a.Before(b) || b.After(a) == false || a.Compare(b) != -1 || a.Equal(b) {
		t.Errorf("valid<valid: a=%v b=%v", a, b)
	}
	if !b.After(a) || a.After(b) || b.Compare(a) != +1 {
		t.Errorf("valid>valid: a=%v b=%v", a, b)
	}
	if !a.Equal(a) || a.Compare(a) != 0 {
		t.Errorf("valid=valid: a=%v", a)
	}

	// NULL vs valid
	if !z.Before(a) || z.After(a) || z.Compare(a) != -1 || z.Equal(a) {
		t.Errorf("NULL<valid: z=%v a=%v", z, a)
	}
	if a.Before(z) || !a.After(z) || a.Compare(z) != +1 || a.Equal(z) {
		t.Errorf("valid>NULL: a=%v z=%v", a, z)
	}

	// NULL vs NULL
	if z.Before(z2) || z.After(z2) || z.Compare(z2) != 0 || !z.Equal(z2) {
		t.Errorf("NULL=NULL: z=%v z2=%v", z, z2)
	}
}

func TestSortableComparison_NullOffsetDateTime(t *testing.T) {
	a := NullOffsetDateTimeFromTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	z := NewNullOffsetDateTimeEmpty()
	if !z.Before(a) || a.Before(z) || z.Compare(a) != -1 || a.Compare(z) != +1 {
		t.Error("sortable NULL semantics broken")
	}
	if !z.Equal(NewNullOffsetDateTimeEmpty()) || z.Compare(NewNullOffsetDateTimeEmpty()) != 0 {
		t.Error("two NULLs must compare equal")
	}
}

func TestSortableComparison_NullOffsetTime(t *testing.T) {
	a := NullOffsetTimeFromTime(time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC))
	z := NewNullOffsetTimeEmpty()
	if !z.Before(a) || a.Before(z) || z.Compare(a) != -1 || a.Compare(z) != +1 {
		t.Error("sortable NULL semantics broken")
	}
	if !z.Equal(NewNullOffsetTimeEmpty()) {
		t.Error("two NULLs must compare equal")
	}
}

func TestSortableComparison_NullLocalDateTime(t *testing.T) {
	a := NewNullLocalDateTime(NewLocalDateTime(2024, time.July, 11, 12, 0, 0, 0))
	z := NewNullLocalDateTimeEmpty()
	if !z.Before(a) || a.Before(z) || z.Compare(a) != -1 || a.Compare(z) != +1 {
		t.Error("sortable NULL semantics broken")
	}
	if !z.Equal(NewNullLocalDateTimeEmpty()) {
		t.Error("two NULLs must compare equal")
	}
}

func TestSortableComparison_NullLocalTime(t *testing.T) {
	a := NewNullLocalTime(NewLocalTime(12, 0, 0, 0))
	z := NewNullLocalTimeEmpty()
	if !z.Before(a) || a.Before(z) || z.Compare(a) != -1 || a.Compare(z) != +1 {
		t.Error("sortable NULL semantics broken")
	}
	if !z.Equal(NewNullLocalTimeEmpty()) {
		t.Error("two NULLs must compare equal")
	}
}

// Comparison smoke tests for not-null types.

func TestComparison_Date(t *testing.T) {
	a := NewDate(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	b := NewDate(time.Date(2024, 7, 11, 0, 0, 0, 0, time.UTC))
	if !a.Before(b) || a.Compare(b) != -1 || a.Equal(b) {
		t.Errorf("Date<Date broken: a=%v b=%v", a, b)
	}
	if !a.Equal(a) || a.Compare(a) != 0 {
		t.Errorf("Date.Equal(self) broken")
	}
}

func TestComparison_OffsetDateTime(t *testing.T) {
	a := NewOffsetDateTime(time.Date(2024, 7, 11, 10, 0, 0, 0, time.FixedZone("UTC+04:00", 4*3600)))
	b := NewOffsetDateTime(time.Date(2024, 7, 11, 6, 0, 0, 0, time.UTC)) // same instant in UTC
	// Equal is instant-based regardless of zone.
	if !a.Equal(b) || a.Compare(b) != 0 {
		t.Errorf("Equal across zones broken: a=%v b=%v", a, b)
	}
}

func TestComparison_LocalDateTime(t *testing.T) {
	a := NewLocalDateTime(2024, time.July, 11, 10, 0, 0, 0)
	b := NewLocalDateTime(2024, time.July, 11, 11, 0, 0, 0)
	if !a.Before(b) || a.Compare(b) != -1 || a.Equal(b) {
		t.Error("LocalDateTime<LocalDateTime broken")
	}
}

func TestComparison_LocalTime(t *testing.T) {
	a := NewLocalTime(10, 0, 0, 0)
	b := NewLocalTime(11, 0, 0, 0)
	if !a.Before(b) || a.Compare(b) != -1 || a.Equal(b) {
		t.Error("LocalTime<LocalTime broken")
	}
}

// Arithmetic round-trip: t.Add(d).Sub(t) == d, AddDate(0,0,0) and
// Add(0) are tautologies, NULL propagates through Add/AddDate/Truncate
// and SubOk reports (0, false).

func TestArithmetic_OffsetDateTime(t *testing.T) {
	src := NewOffsetDateTime(time.Date(2024, 7, 11, 10, 0, 0, 0, time.UTC))
	d := 3 * time.Hour
	if got := src.Add(d).Sub(src); got != d {
		t.Errorf("Add/Sub round-trip: expected %v, got %v", d, got)
	}
	if !src.AddDate(0, 0, 0).Equal(src) {
		t.Error("AddDate(0,0,0) must be tautology")
	}
	if !src.Add(0).Equal(src) {
		t.Error("Add(0) must be tautology")
	}

	// Truncate to hour drops minutes/seconds/ns.
	with := NewOffsetDateTime(time.Date(2024, 7, 11, 10, 37, 42, 999, time.UTC))
	tr := with.Truncate(time.Hour)
	if !tr.Equal(NewOffsetDateTime(time.Date(2024, 7, 11, 10, 0, 0, 0, time.UTC))) {
		t.Errorf("Truncate(1h) broken: got %v", tr.ToString())
	}
}

func TestArithmetic_NullOffsetDateTime_NULLPropagation(t *testing.T) {
	z := NewNullOffsetDateTimeEmpty()
	if z.Add(time.Hour).Valid {
		t.Error("Add on NULL must stay NULL")
	}
	if z.AddDate(0, 0, 1).Valid {
		t.Error("AddDate on NULL must stay NULL")
	}
	if z.Truncate(time.Hour).Valid {
		t.Error("Truncate on NULL must stay NULL")
	}

	a := NullOffsetDateTimeFromTime(time.Now())
	if d, ok := z.SubOk(a); ok || d != 0 {
		t.Errorf("SubOk(NULL, valid) must be (0, false), got (%v, %v)", d, ok)
	}
	if d, ok := a.SubOk(z); ok || d != 0 {
		t.Errorf("SubOk(valid, NULL) must be (0, false), got (%v, %v)", d, ok)
	}
	if d, ok := a.SubOk(a); !ok || d != 0 {
		t.Errorf("SubOk(self, self) must be (0, true), got (%v, %v)", d, ok)
	}
}

func TestArithmetic_Date(t *testing.T) {
	d1 := NewDate(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	d2 := d1.AddDate(0, 0, 7)
	if got := d2.Sub(d1) / (24 * time.Hour); got != 7 {
		t.Errorf("Date.AddDate/Sub: expected 7 days, got %d", got)
	}
}

func TestArithmetic_LocalDateTime(t *testing.T) {
	src := NewLocalDateTime(2024, time.July, 11, 10, 0, 0, 0)
	d := 3 * time.Hour
	if got := src.Add(d).Sub(src); got != d {
		t.Errorf("LocalDateTime Add/Sub round-trip: %v", got)
	}
}

func TestArithmetic_TimeOnly_WrapDocumented(t *testing.T) {
	// OffsetTime.Add wraps past midnight without enforcement — verify
	// the result is well-defined (matches time.Time.Add) rather than
	// rejected. The serialised form drops the date, so a 25h offset
	// shows as the H:M:S of the next day.
	src := NewOffsetTime(time.Date(0, 1, 1, 1, 0, 0, 0, time.UTC))
	got := src.Add(25 * time.Hour)
	want := NewOffsetTime(time.Date(0, 1, 2, 2, 0, 0, 0, time.UTC))
	if !got.Equal(want) {
		t.Errorf("OffsetTime.Add(25h): expected %v, got %v", want.ToString(), got.ToString())
	}
}

// LocalTime arithmetic: Add/Sub round-trip, Truncate.
func TestArithmetic_LocalTime(t *testing.T) {
	src := NewLocalTime(10, 0, 0, 0)
	if got := src.Add(2 * time.Hour).Sub(src); got != 2*time.Hour {
		t.Errorf("LocalTime Add/Sub: %v", got)
	}
	if got := src.Truncate(time.Hour); got != src {
		t.Errorf("LocalTime Truncate(1h) on whole-hour: %v", got)
	}
}

func TestArithmetic_NullDate_Truncate_NULL(t *testing.T) {
	z := NewNullDateEmpty()
	a := NullDateFromTime(time.Date(2024, 7, 11, 0, 0, 0, 0, time.UTC))
	if z.Add(time.Hour).Valid {
		t.Error("Add on NULL must stay NULL")
	}
	if z.AddDate(0, 0, 1).Valid {
		t.Error("AddDate on NULL must stay NULL")
	}
	if d, ok := z.SubOk(a); ok || d != 0 {
		t.Errorf("SubOk(NULL, valid) must be (0, false), got (%v, %v)", d, ok)
	}
}

func TestArithmetic_NullLocalTime_NULL(t *testing.T) {
	z := NewNullLocalTimeEmpty()
	if z.Add(time.Hour).Valid {
		t.Error("Add on NULL must stay NULL")
	}
	if z.Truncate(time.Hour).Valid {
		t.Error("Truncate on NULL must stay NULL")
	}
}

func TestArithmetic_NullLocalDateTime_NULL(t *testing.T) {
	z := NewNullLocalDateTimeEmpty()
	if z.Add(time.Hour).Valid || z.AddDate(0, 0, 1).Valid || z.Truncate(time.Hour).Valid {
		t.Error("NULL passthrough broken")
	}
	a := NewNullLocalDateTime(NewLocalDateTime(2024, time.July, 11, 12, 0, 0, 0))
	if d, ok := z.SubOk(a); ok || d != 0 {
		t.Errorf("SubOk on NULL operand: (%v, %v)", d, ok)
	}
}

func TestArithmetic_NullOffsetTime_NULL(t *testing.T) {
	z := NewNullOffsetTimeEmpty()
	if z.Add(time.Hour).Valid || z.Truncate(time.Hour).Valid {
		t.Error("NULL passthrough broken")
	}
	a := NullOffsetTimeFromTime(time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC))
	if d, ok := z.SubOk(a); ok || d != 0 {
		t.Errorf("SubOk on NULL operand: (%v, %v)", d, ok)
	}
}

// Unix family on not-null types.
func TestUnix_NotNull(t *testing.T) {
	odt := NewOffsetDateTime(time.Date(2024, 1, 1, 0, 0, 0, 500_000, time.UTC))
	if odt.Unix() != 1704067200 {
		t.Errorf("OffsetDateTime.Unix: %d", odt.Unix())
	}
	if odt.UnixMilli() != 1704067200_000 {
		t.Errorf("OffsetDateTime.UnixMilli: %d", odt.UnixMilli())
	}
	if odt.UnixMicro() != 1704067200_000_500 {
		t.Errorf("OffsetDateTime.UnixMicro: %d", odt.UnixMicro())
	}
	if odt.UnixNano() != 1704067200_000_500_000 {
		t.Errorf("OffsetDateTime.UnixNano: %d", odt.UnixNano())
	}

	d := NewDate(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	if d.Unix() != 1704067200 || d.UnixMilli() != 1704067200_000 ||
		d.UnixMicro() != 1704067200_000_000 || d.UnixNano() != 1704067200_000_000_000 {
		t.Error("Date Unix family broken")
	}
}

// Unix family on null types — Ok-variants distinguish NULL from epoch.
func TestUnix_NullOk(t *testing.T) {
	z := NewNullDateEmpty()
	if v, ok := z.UnixOk(); ok || v != 0 {
		t.Errorf("NULL UnixOk: (%d, %v)", v, ok)
	}
	if v, ok := z.UnixMilliOk(); ok || v != 0 {
		t.Errorf("NULL UnixMilliOk: (%d, %v)", v, ok)
	}
	if v, ok := z.UnixMicroOk(); ok || v != 0 {
		t.Errorf("NULL UnixMicroOk: (%d, %v)", v, ok)
	}
	if v, ok := z.UnixNanoOk(); ok || v != 0 {
		t.Errorf("NULL UnixNanoOk: (%d, %v)", v, ok)
	}

	epoch := NullDateFromTime(time.Unix(0, 0).UTC())
	if v, ok := epoch.UnixOk(); !ok || v != 0 {
		t.Errorf("epoch UnixOk: (%d, %v) — must distinguish NULL from epoch", v, ok)
	}

	zod := NewNullOffsetDateTimeEmpty()
	if v, ok := zod.UnixMilliOk(); ok || v != 0 {
		t.Errorf("NULL OffsetDateTime UnixMilliOk: (%d, %v)", v, ok)
	}
	if v, ok := zod.UnixMicroOk(); ok || v != 0 {
		t.Errorf("NULL OffsetDateTime UnixMicroOk: (%d, %v)", v, ok)
	}
	if v, ok := zod.UnixNanoOk(); ok || v != 0 {
		t.Errorf("NULL OffsetDateTime UnixNanoOk: (%d, %v)", v, ok)
	}
}

// Component getters smoke test.
func TestGetters(t *testing.T) {
	odt := NewOffsetDateTime(time.Date(2024, time.July, 11, 13, 7, 42, 999, time.UTC))
	if odt.Year() != 2024 || odt.Month() != time.July || odt.Day() != 11 ||
		odt.Hour() != 13 || odt.Minute() != 7 || odt.Second() != 42 ||
		odt.Nanosecond() != 999 || odt.YearDay() != 193 {
		t.Errorf("OffsetDateTime getters broken: %+v", odt)
	}

	ot := NewOffsetTime(time.Date(0, 1, 1, 13, 7, 42, 999, time.UTC))
	if ot.Hour() != 13 || ot.Minute() != 7 || ot.Second() != 42 || ot.Nanosecond() != 999 {
		t.Errorf("OffsetTime getters broken: %+v", ot)
	}

	d := NewDate(time.Date(2024, time.July, 11, 0, 0, 0, 0, time.UTC))
	if d.Year() != 2024 || d.Month() != time.July || d.Day() != 11 || d.YearDay() != 193 {
		t.Errorf("Date getters broken: %+v", d)
	}

	ldt := NewLocalDateTime(2024, time.July, 11, 13, 7, 42, 0)
	if ldt.YearDay() != 193 {
		t.Errorf("LocalDateTime.YearDay: %d", ldt.YearDay())
	}
}

// Equal cross-matrix on null types — fills the gap left by smoke
// tests above (which only check 2 of 4 cases). Each case asserts
// Equal returns the expected bool.
func TestEqualMatrix_Nulls(t *testing.T) {
	// NullOffsetDateTime
	a := NullOffsetDateTimeFromTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	b := NullOffsetDateTimeFromTime(time.Date(2024, 7, 11, 0, 0, 0, 0, time.UTC))
	z := NewNullOffsetDateTimeEmpty()
	if !a.Equal(a) || a.Equal(b) || a.Equal(z) || z.Equal(a) {
		t.Error("NullOffsetDateTime.Equal matrix broken")
	}

	// NullOffsetTime
	at := NullOffsetTimeFromTime(time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC))
	bt := NullOffsetTimeFromTime(time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC))
	zt := NewNullOffsetTimeEmpty()
	if !at.Equal(at) || at.Equal(bt) || at.Equal(zt) || zt.Equal(at) {
		t.Error("NullOffsetTime.Equal matrix broken")
	}

	// NullLocalDateTime
	al := NewNullLocalDateTime(NewLocalDateTime(2024, time.July, 11, 12, 0, 0, 0))
	bl := NewNullLocalDateTime(NewLocalDateTime(2024, time.July, 11, 13, 0, 0, 0))
	zl := NewNullLocalDateTimeEmpty()
	if !al.Equal(al) || al.Equal(bl) || al.Equal(zl) || zl.Equal(al) {
		t.Error("NullLocalDateTime.Equal matrix broken")
	}

	// NullLocalTime
	alt := NewNullLocalTime(NewLocalTime(12, 0, 0, 0))
	blt := NewNullLocalTime(NewLocalTime(13, 0, 0, 0))
	zlt := NewNullLocalTimeEmpty()
	if !alt.Equal(alt) || alt.Equal(blt) || alt.Equal(zlt) || zlt.Equal(alt) {
		t.Error("NullLocalTime.Equal matrix broken")
	}
}

// Compare strict-greater branch on null types — fills the gap where
// only the strict-less branch was exercised.
func TestCompare_StrictGreater_Nulls(t *testing.T) {
	a := NullOffsetTimeFromTime(time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC))
	b := NullOffsetTimeFromTime(time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC))
	if b.Compare(a) != +1 || !b.After(a) {
		t.Error("NullOffsetTime: valid>valid broken")
	}

	al := NewNullLocalDateTime(NewLocalDateTime(2024, time.July, 11, 12, 0, 0, 0))
	bl := NewNullLocalDateTime(NewLocalDateTime(2024, time.July, 11, 13, 0, 0, 0))
	if bl.Compare(al) != +1 || !bl.After(al) {
		t.Error("NullLocalDateTime: valid>valid broken")
	}

	alt := NewNullLocalTime(NewLocalTime(12, 0, 0, 0))
	blt := NewNullLocalTime(NewLocalTime(13, 0, 0, 0))
	if blt.Compare(alt) != +1 || !blt.After(alt) {
		t.Error("NullLocalTime: valid>valid broken")
	}
}

// Valid-receiver arithmetic on null types — fills the gap where only
// NULL-receiver branches were exercised.
func TestArithmetic_NullValidReceiver(t *testing.T) {
	// NullDate
	d := NullDateFromTime(time.Date(2024, 7, 11, 0, 0, 0, 0, time.UTC))
	if got := d.Add(24 * time.Hour); !got.Valid ||
		!got.Val.Equal(NewDate(time.Date(2024, 7, 12, 0, 0, 0, 0, time.UTC))) {
		t.Errorf("NullDate.Add valid: %v", got.ToString())
	}
	if got := d.AddDate(0, 0, 1); !got.Valid ||
		!got.Val.Equal(NewDate(time.Date(2024, 7, 12, 0, 0, 0, 0, time.UTC))) {
		t.Errorf("NullDate.AddDate valid: %v", got.ToString())
	}
	d2 := NullDateFromTime(time.Date(2024, 7, 18, 0, 0, 0, 0, time.UTC))
	if dur, ok := d2.SubOk(d); !ok || dur != 7*24*time.Hour {
		t.Errorf("NullDate.SubOk valid: (%v, %v)", dur, ok)
	}
	if v, ok := d.UnixOk(); !ok || v != 1720656000 {
		t.Errorf("NullDate.UnixOk valid: (%d, %v)", v, ok)
	}
	if v, ok := d.UnixMilliOk(); !ok || v != 1720656000_000 {
		t.Errorf("NullDate.UnixMilliOk valid: (%d, %v)", v, ok)
	}
	if v, ok := d.UnixMicroOk(); !ok || v != 1720656000_000_000 {
		t.Errorf("NullDate.UnixMicroOk valid: (%d, %v)", v, ok)
	}
	if v, ok := d.UnixNanoOk(); !ok || v != 1720656000_000_000_000 {
		t.Errorf("NullDate.UnixNanoOk valid: (%d, %v)", v, ok)
	}

	// NullOffsetDateTime
	od := NullOffsetDateTimeFromTime(time.Date(2024, 7, 11, 10, 0, 0, 0, time.UTC))
	if got := od.Add(time.Hour); !got.Valid {
		t.Error("NullOffsetDateTime.Add valid stays NULL")
	}
	if got := od.AddDate(0, 0, 1); !got.Valid {
		t.Error("NullOffsetDateTime.AddDate valid stays NULL")
	}
	if got := od.Truncate(time.Hour); !got.Valid {
		t.Error("NullOffsetDateTime.Truncate valid stays NULL")
	}
	if v, ok := od.UnixMilliOk(); !ok || v == 0 {
		t.Errorf("NullOffsetDateTime.UnixMilliOk valid: (%d, %v)", v, ok)
	}
	if v, ok := od.UnixMicroOk(); !ok || v == 0 {
		t.Errorf("NullOffsetDateTime.UnixMicroOk valid: (%d, %v)", v, ok)
	}
	if v, ok := od.UnixNanoOk(); !ok || v == 0 {
		t.Errorf("NullOffsetDateTime.UnixNanoOk valid: (%d, %v)", v, ok)
	}

	// NullOffsetTime
	ot := NullOffsetTimeFromTime(time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC))
	if got := ot.Add(time.Hour); !got.Valid {
		t.Error("NullOffsetTime.Add valid stays NULL")
	}
	if got := ot.Truncate(time.Hour); !got.Valid {
		t.Error("NullOffsetTime.Truncate valid stays NULL")
	}
	ot2 := NullOffsetTimeFromTime(time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC))
	if dur, ok := ot2.SubOk(ot); !ok || dur != time.Hour {
		t.Errorf("NullOffsetTime.SubOk valid: (%v, %v)", dur, ok)
	}

	// NullLocalDateTime
	ldt := NewNullLocalDateTime(NewLocalDateTime(2024, time.July, 11, 12, 0, 0, 0))
	if got := ldt.Add(time.Hour); !got.Valid {
		t.Error("NullLocalDateTime.Add valid stays NULL")
	}
	if got := ldt.AddDate(0, 0, 1); !got.Valid {
		t.Error("NullLocalDateTime.AddDate valid stays NULL")
	}
	if got := ldt.Truncate(time.Hour); !got.Valid {
		t.Error("NullLocalDateTime.Truncate valid stays NULL")
	}
	ldt2 := NewNullLocalDateTime(NewLocalDateTime(2024, time.July, 11, 13, 0, 0, 0))
	if dur, ok := ldt2.SubOk(ldt); !ok || dur != time.Hour {
		t.Errorf("NullLocalDateTime.SubOk valid: (%v, %v)", dur, ok)
	}

	// NullLocalTime
	lt := NewNullLocalTime(NewLocalTime(12, 0, 0, 0))
	if got := lt.Add(time.Hour); !got.Valid {
		t.Error("NullLocalTime.Add valid stays NULL")
	}
	if got := lt.Truncate(time.Hour); !got.Valid {
		t.Error("NullLocalTime.Truncate valid stays NULL")
	}
}

// UTC() shorthand.
func TestUTC_Shorthand(t *testing.T) {
	plus4 := time.FixedZone("UTC+04:00", 4*3600)
	odt := NewOffsetDateTime(time.Date(2024, 7, 11, 10, 0, 0, 0, plus4))
	if got := odt.UTC().ToString(); got != "2024-07-11T06:00:00Z" {
		t.Errorf("OffsetDateTime.UTC: %q", got)
	}

	ot := NewOffsetTime(time.Date(0, 1, 1, 10, 0, 0, 0, plus4))
	if got := ot.UTC().ToString(); got != "06:00:00Z" {
		t.Errorf("OffsetTime.UTC: %q", got)
	}

	// Null variants.
	zod := NewNullOffsetDateTimeEmpty()
	if zod.UTC().Valid {
		t.Error("NULL UTC must stay NULL")
	}
	zot := NewNullOffsetTimeEmpty()
	if zot.UTC().Valid {
		t.Error("NULL UTC must stay NULL")
	}
}
