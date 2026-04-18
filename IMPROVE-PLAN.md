# `gxsrvs/dtx` — Improvement Plan

This document lists the work required to take the library to a public
`v1.0.0` release. It builds on `REVIEW.md` (current-state assessment) and
`TESTS-PLAN.md` (test coverage roadmap).

The plan is organised by topic (**functional**, **architectural**,
**infrastructural**) and split into milestones M1..M14.

Priorities:

- 🔴 blocker for the public release,
- 🟡 important, but can land shortly after,
- 🟢 nice-to-have.

---

## M1. 🔴 Non-null `Date`, `Time`, `DateTime` using the library’s serializer

**Motivation.** Today a consumer has to use bare `time.Time` for a required
datetime field (which marshals via stdlib RFC3339, carrying a timezone),
while optional fields use `types.NullDate` / `types.NullDateTime` /
`types.NullTime` with the library’s own format. The two sides disagree,
and the client has to parse both shapes. A non-null type that shares the
same serializer with its nullable counterpart eliminates the split.

**API sketch.**

```go
// types/date.go
type Date time.Time

func NewDate(t time.Time) Date
func DateFromString(s string) (Date, error)        // same formats as NullDate
func (Date) MarshalJSON() ([]byte, error)          // "2006-01-02"
func (*Date) UnmarshalJSON([]byte) error
func (Date) Value() (driver.Value, error)          // time.Time
func (*Date) Scan(interface{}) error
func (Date) ToString() string
func (Date) AsTime() time.Time

// types/time.go — time-of-day
type Time time.Time
// same formats as NullTime (HH:MM[:SS[.ms]][+HH:MM])

// types/datetime.go — update the existing DateTime:
//   — add Value(),
//   — remove the dead if/else branch in Scan,
//   — align the format with NullDateTime.
```

**Principle.** `Date.MarshalJSON` and `NullDate.MarshalJSON` share a single
private helper `formatDate(time.Time) string`. Likewise for `Time`,
`DateTime`, `LocalDateTime`. Parsing goes through `parseDate`, `parseTime`,
`parseDateTime`, `parseLocalDateTime` — one source of truth each.

**Subtle point.** Unlike the nullable sibling, a non-null type cannot be
“empty”. `UnmarshalJSON([]byte("null"))` must return an error. Today
`DateTime.UnmarshalJSON` silently accepts `null` as a no-op; M1 fixes that.

---

## M2. 🔴 `LocalDateTime` and `NullLocalDateTime` (no timezone)

**Motivation.** Business events are often anchored to a local wall-clock
context (a contract is signed “on 2026-04-18 at 13:00” — a moment in the
wall-clock time of an office, with no offset attached). `time.Time`
always carries a zone and is the wrong model.

**Internal representation.**

Option A — a `time.Time` wrapper forced to `time.UTC` that silently drops
the zone on marshal. The risk is that any comparison with an ordinary
`time.Time` elsewhere may yield a wrong answer.

Option B (preferred) — a dedicated struct:

```go
type LocalDateTime struct {
    Year    int
    Month   time.Month
    Day     int
    Hour    int
    Minute  int
    Second  int
    Nanosec int
}
```

Pros: explicit semantics, no zone leaks through `time.Time`, guards against
accidentally using `time.Now()` as a value source. Con: slightly more code.

**API sketch.**

```go
func NewLocalDateTime(y int, m time.Month, d, hh, mm, ss, ns int) LocalDateTime
func LocalDateTimeFromString(s string) (LocalDateTime, error)  // "2006-01-02T15:04:05" / "2006-01-02 15:04:05"
func (LocalDateTime) ToTime(loc *time.Location) time.Time
func (LocalDateTime) MarshalJSON() ([]byte, error)             // "2006-01-02T15:04:05"
func (*LocalDateTime) UnmarshalJSON([]byte) error
func (LocalDateTime) Value() (driver.Value, error)
func (*LocalDateTime) Scan(interface{}) error
func (LocalDateTime) ToString() string

type NullLocalDateTime struct {
    Val   LocalDateTime
    Valid bool
}

// standard set: New, NewEmpty, FromString, IsEmpty, ToString,
// Value, Scan, MarshalJSON, UnmarshalJSON.
```

**Format.** `YYYY-MM-DDTHH:MM:SS[.fff]` with no offset, and no `Z` /
`+NN:NN` suffix. A string carrying a TZ is rejected with an error to keep
the semantics explicit.

---

## M3. 🔴 Fix `Scan` in every nullable type

Replace `reflect.TypeOf(value) == nil` with the real source of truth —
the `Valid` field of the embedded `sql.Null*`. For `NullDate`:

```go
func (v *NullDate) Scan(value interface{}) error {
    var s sql.NullTime
    if err := s.Scan(value); err != nil {
        return err
    }
    if !s.Valid {
        *v = NewNullDateEmpty()
        return nil
    }
    *v = NewNullDate(s.Time)
    return nil
}
```

Affected files: `null_date.go`, `null_datetime.go`, `null_time.go`,
`null_iso_date.go`, `null_string.go`, `null_bool.go`,
`null_int{16,32,64}.go`, `null_float.go`, `null_decimal.go`,
`null_uuid.go`.

Cover with tests (see `TESTS-PLAN.md` §4.1).

---

## M4. 🟡 Consistent `ToString()`

Rules:

- `Valid == false` → `""` across every type (today `NullFloat` returns
  `"null"`),
- for date/time types — the **same** format as `MarshalJSON`
  (today `NullTime.ToString` calls `time.Time.String()`, while JSON uses
  `Format(time.TimeOnly)`).

Affected: `null_float.go`, `null_time.go`.

---

## M5. 🟡 Rename `NullTime` for semantic clarity

In this library `NullTime` is time-of-day; in `database/sql` `NullTime`
means timestamp. The confusion is baked into the name. Options:

1. Rename `NullTime` → `NullTimeOfDay` / `NullClockTime`.
2. Keep the name, but document the meaning clearly in godoc and README.

Because the library has no external consumers yet, option 1 is preferred.
Until v1.0 this is an acceptable breaking change.

Similarly, an offset-less `NullDateTime` is misleading. M2 resolves that
by introducing `NullLocalDateTime`.

---

## M6. 🟡 Remove `//goland:noinspection` and stabilise receivers

One receiver style per type. Recommendation:

- **value receiver** for `Marshal*`, `Value()`, `ToString()` — read-only
  API,
- **pointer receiver** for `Scan`, `Unmarshal*`, `IsEmpty` — mutating,
  and cheaper copies on larger types.

Done once; afterwards every `//goland:noinspection GoMixedReceiverTypes`
disappears.

---

## M7. 🟡 Interface dispatch for `IsEmpty` / `ToString`

Drop the hand-rolled type switch:

```go
func IsEmpty(v interface{}) bool {
    if v == nil {
        return true
    }
    if e, ok := v.(Emptiable); ok {
        return e.IsEmpty()
    }
    // fallback: stdlib nullables, primitives
    ...
}
```

Same shape for `ToString`. Adding a new type no longer requires edits to
`emptiable.go` / `string_able.go`.

---

## M8. 🟡 `utils.ToJson` should return an error

Current signature:

```go
func ToJson(entity interface{}) string
```

Swallowing the error is an anti-pattern for a public API. Replace with:

```go
func ToJson(entity interface{}) (string, error)
// and, optionally, MustToJson(entity interface{}) string for test helpers
```

---

## M9. 🟡 Add `LICENSE` + overhaul `README.md`

Add a `LICENSE` file (MIT). Update the README per `REVIEW.md` §6:

- badges (`pkg.go.dev`, Go Report Card, CI status, coverage),
- full documentation of **every** public type,
- sections Installation / Quickstart / Usage / SQL / JSON / TZ / Testing,
- semver policy,
- contributing,
- English as the sole language of the docs.

---

## M10. 🟡 godoc coverage

Every public function / type / constant gets a `// Name ...` godoc
comment. Add `// Package types provides nullable wrappers ...` in a new
`types/doc.go`, and do the same in `dto/` and `utils/`.

---

## M11. 🟡 CI/CD

`.github/workflows/ci.yml`:

- test matrix (Go 1.22..1.26),
- `golangci-lint` (gofmt, govet, staticcheck, errcheck, unused, goimports,
  revive),
- Codecov upload,
- release workflow (goreleaser or `release-please`).

Adjust `go.mod` — drop `go 1.26.1` to the lowest supported minor
(suggest `go 1.22`).

---

## M12. 🟢 Versioning

- Switch to semver: `v0.x` tags until the API stabilises, `v1.0.0`
  once M1..M9 are done.
- Introduce `CHANGELOG.md` in the “Keep a Changelog” format.
- Each PR updates the `Unreleased` section.

---

## M13. 🟢 API surface expansion

Driven by real-world usage:

- `NullByte` / `NullRune` — if needed.
- `NullJSON[T any]` — generic wrapper for fields that store a serialised
  entity.
- `Equal(other)` on every nullable.
- `Compare(a, b)` — for sorting.
- `NullDate.AddDays(n int)`, `NullDateTime.AddDuration`.
- `Date.Today()`, `DateTime.Now()`, `LocalDateTime.Now(loc)`.

---

## M14. 🟢 `dto.DataPackage` evolution

- Drop `DataObject = interface{}` — it contributes nothing. Keep
  `DataPackage[T any]`.
- Add standard envelopes for REST responses: `Result[T]`,
  `PaginatedResult[T]` (with `Total`, `Limit`, `Offset`), `ErrorResponse`.
- JSON round-trip tests with real nullable types.

---

## Roadmap

| Milestone | Description                                                  | Priority | Size |
| --------- | ------------------------------------------------------------ | :------: | :--: |
| M1        | Non-null Date/Time/DateTime + shared serializer              |    🔴    |  M   |
| M2        | LocalDateTime / NullLocalDateTime                            |    🔴    |  L   |
| M3        | `Scan` fix (use `sql.Null*.Valid`)                           |    🔴    |  M   |
| M4        | Consistent `ToString`                                        |    🟡    |  S   |
| M5        | `NullTime` rename (pre-v1.0 breaking change)                 |    🟡    |  S   |
| M6        | Receivers + remove `//goland:noinspection`                   |    🟡    |  S   |
| M7        | Interface dispatch for `IsEmpty` / `ToString`                |    🟡    |  S   |
| M8        | `utils.ToJson` returns an error                              |    🟡    |  S   |
| M9        | LICENSE + README                                             |    🟡    |  S   |
| M10       | godoc across the public API                                  |    🟡    |  M   |
| M11       | CI/CD + linter + coverage                                    |    🟡    |  M   |
| M12       | Semver + CHANGELOG                                           |    🟢    |  S   |
| M13       | API expansion (Equal / Compare / arithmetic)                 |    🟢    |  M   |
| M14       | `dto` evolution (Result, PaginatedResult)                    |    🟢    |  M   |

Sizes: S ≈ up to one day, M ≈ two–three days, L ≈ up to a week.

**Minimum work to reach the public release:** M1..M3 + M9 + M11, and ideally
M10. Everything else can land iteratively.