// Package bench measures the performance of dtx Null* wrappers
// against the idiomatic Go alternative — using pointers (*T) for
// nullable struct fields. Each domain entity appears in two variants
// (one dtx-typed, one pointer-typed) and every benchmark runs over
// three input profiles (all-valid, all-null, mixed) so the comparison
// reflects realistic fill ratios rather than a single shape.
//
// The pointer variant uses the conventions a Go developer would
// actually reach for in production: omitempty on every nullable field
// and time.Time for date/datetime fields. As a consequence the JSON
// shapes are not byte-identical between the two variants — pointer
// JSON omits absent fields, dtx JSON emits them as explicit "null",
// and time.Time always serialises as RFC 3339 rather than the
// date-only layout dtx uses for NullDate. The TestClientJSONSizes
// helper reports the resulting byte sizes so the benchmark numbers
// can be read in context.
//
// Run with:
//
//	go test ./bench/... -bench=. -benchmem
//
// To filter to a single variant, profile, or operation, pass a more
// specific -bench expression:
//
//	go test ./bench/... -bench=Client/Null/Unmarshal/Mixed -benchmem
package bench
