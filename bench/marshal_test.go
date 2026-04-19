package bench

import (
	"testing"
	"time"

	"github.com/gxsrvs/dtx/types"
)

// Focused micro-benchmarks for individual Null* MarshalJSON
// implementations. Use these to validate per-type optimisations in
// isolation before re-running the higher-level Client benchmark.
//
//	go test ./bench/ -bench=BenchmarkMarshal -benchmem -run=^$

var sampleNullInt16Valid = types.NewNullInt16(1969)
var sampleNullInt32Valid = types.NewNullInt32(1969_07_21)
var sampleNullInt64Valid = types.NewNullInt64(1969_07_21)
var sampleNullInt64Empty = types.NewNullInt64Empty()

var sampleNullBoolValid = types.NewNullBool(true)
var sampleNullFloatValid = types.NewNullFloat(3.141592653589793)

var sampleNullDateValid = types.NewNullDate(
	time.Date(1961, 4, 12, 0, 0, 0, 0, time.UTC),
)
var sampleNullDateEmpty = types.NewNullDateEmpty()

var sampleNullLocalDateTimeValid = types.NewNullLocalDateTime(
	types.LocalDateTime{Year: 1969, Month: 7, Day: 21, Hour: 2, Minute: 56, Second: 15},
)
var sampleNullLocalTimeValid = types.NewNullLocalTime(
	types.LocalTime{Hour: 2, Minute: 56, Second: 15},
)
var sampleNullOffsetDateTimeValid = types.NewNullOffsetDateTime(
	time.Date(1969, 7, 21, 2, 56, 15, 0, time.UTC),
)
var sampleNullOffsetTimeValid = types.NewNullOffsetTime(
	time.Date(0, 1, 1, 2, 56, 15, 0, time.UTC),
)

func BenchmarkMarshalNullInt64Valid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := sampleNullInt64Valid.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalNullInt64Empty(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := sampleNullInt64Empty.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalNullDateValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := sampleNullDateValid.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalNullDateEmpty(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := sampleNullDateEmpty.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalNullInt16Valid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := sampleNullInt16Valid.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalNullInt32Valid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := sampleNullInt32Valid.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalNullBoolValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := sampleNullBoolValid.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalNullFloatValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := sampleNullFloatValid.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalNullLocalDateTimeValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := sampleNullLocalDateTimeValid.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalNullLocalTimeValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := sampleNullLocalTimeValid.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalNullOffsetDateTimeValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := sampleNullOffsetDateTimeValid.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalNullOffsetTimeValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := sampleNullOffsetTimeValid.MarshalJSON(); err != nil {
			b.Fatal(err)
		}
	}
}
