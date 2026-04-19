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

## Canonical type set (locked)

The following layout is now the canonical naming for date/time types:

| Not-null            | Nullable                | TZ?  | JSON format                          |
| ------------------- | ----------------------- | :--: | ------------------------------------ |
| `Date`              | `NullDate`              |  —   | `2006-01-02`                         |
| `OffsetTime`        | `NullOffsetTime`        | yes  | `15:04:05[.fff…]Z` / `±HH:MM`        |
| `OffsetDateTime`    | `NullOffsetDateTime`    | yes  | RFC 3339 (`2006-01-02T15:04:05…Z`)   |
| `LocalTime`         | `NullLocalTime`         |  no  | `15:04:05[.fff…]`                    |
| `LocalDateTime`     | `NullLocalDateTime`     |  no  | `2006-01-02T15:04:05[.fff…]`         |

Each not-null type shares its serializer with the nullable counterpart via
private `formatXxx` / `parseXxx` helpers. `Local*` parsers reject inputs
that carry a TZ designator.

---

## M1. ✅ Non-null `Date`, `OffsetTime`, `OffsetDateTime` with shared serializer

Done. `Date` (date.go), `OffsetTime` (offset_time.go), `OffsetDateTime`
(offset_datetime.go) all share their serializer with the nullable
counterpart through `format_date.go`, `format_offset_time.go`,
`format_offset_datetime.go`. Not-null `UnmarshalJSON("null")` returns an
error.

---

## M2. ✅ `LocalDateTime` / `NullLocalDateTime` and `LocalTime` / `NullLocalTime`

Done. Both pairs are backed by a dedicated struct (no `time.Time` inside,
no zone leak):

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

type LocalTime struct {
    Hour    int
    Minute  int
    Second  int
    Nanosec int
}
```

Format is `YYYY-MM-DDTHH:MM:SS[.fff…]` (and `HH:MM:SS[.fff…]` for
`LocalTime`). Inputs with a TZ designator are rejected; `hasOffsetSuffix`
in `types_utils.go` is the shared detector.

---

## M3. ✅ Fix `Scan` in every nullable type

Done. Every nullable `Scan` now derives `Valid` from the embedded
`sql.Null*.Valid` field rather than `reflect.TypeOf(value) == nil`. For
`NullDate`:

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

Files updated: `null_date.go`, `null_iso_date.go`, `null_string.go`,
`null_bool.go`, `null_int{16,32,64}.go`, `null_float.go`,
`null_decimal.go`, `null_uuid.go`. Covered by `scan_reset_test.go` plus
the per-type Scan tests.

---

## M4. 🟡 Consistent `ToString()`

Rules:

- `Valid == false` → `""` across every type (today `NullFloat` returns
  `"null"`),
- for date/time types — the **same** format as `MarshalJSON`.

Status: ✅ done for the date/time family (`NullOffsetTime` now shares
`formatOffsetTime` with `MarshalJSON`).

Still 🟡: `null_float.go` returns `"null"` for invalid values where every
other type returns `""`.

---

## M5. ✅ Rename ambiguous date/time types

Done:

- `TimeOnly` → `OffsetTime`
- `NullTime` → `NullOffsetTime`
- `DateTime` → `OffsetDateTime`
- `NullDateTime` → `NullOffsetDateTime`

The chosen `Offset*` prefix makes the presence of a TZ offset explicit and
avoids the `database/sql.NullTime` (timestamp) vs library-`NullTime`
(time-of-day) confusion.

---

## M6. 🟡 Remove `//goland:noinspection` and stabilise receivers

One receiver style per type:

- **value receiver** for `Marshal*`, `Value()`, `ToString()` — read-only
  API,
- **pointer receiver** for `Scan`, `Unmarshal*`, `IsEmpty` — mutating,
  and cheaper copies on larger types.

Status: ✅ for the new files (`date.go`, `offset_time.go`,
`null_offset_time.go`, `offset_datetime.go`, `null_offset_datetime.go`,
`local_*`).

Still 🟡: the older nullable types (`null_string.go`, `null_bool.go`,
`null_int*.go`, `null_float.go`, `null_decimal.go`, `null_uuid.go`,
`null_iso_date.go`, `null_date.go`) still carry
`//goland:noinspection GoMixedReceiverTypes` and inconsistent receivers.

---

## M7. ✅ Interface dispatch for `IsEmpty` / `ToString`

Done. `IsEmpty` and `ToString` now dispatch through the `Emptiable` /
`ToStringAble` interfaces. The previous hand-rolled `type switch` over
the catalogue panicked at runtime for the value branches of
pointer-receiver types (e.g. `IsEmpty(NullString{})`); the new
implementation falls back to wrapping the value in a fresh pointer via
`reflect.New` so pointer-receiver methods are reachable.

Adding a new type no longer requires edits to `emptiable.go` /
`string_able.go`.

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

`LICENSE` (MIT) and `THIRD_PARTY_LICENSES.md` are already in place. The
remaining work is the README overhaul per `REVIEW.md` §6:

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
| M1        | Non-null Date/OffsetTime/OffsetDateTime + shared serializer  |    ✅    |  M   |
| M2        | LocalDateTime/NullLocalDateTime + LocalTime/NullLocalTime    |    ✅    |  L   |
| M3        | `Scan` fix (use `sql.Null*.Valid`) across every nullable     |    ✅    |  M   |
| M4        | Consistent `ToString` — date/time done, `NullFloat` 🟡       |    🟡    |  S   |
| M5        | Rename `Time*` / `DateTime*` → `Offset*` (breaking, pre-v1.0)|    ✅    |  S   |
| M6        | Receivers + remove `//goland:noinspection` — new files done  |    🟡    |  S   |
| M7        | Interface dispatch for `IsEmpty` / `ToString`                |    ✅    |  S   |
| M8        | `utils.ToJson` returns an error                              |    🟡    |  S   |
| M9        | LICENSE + README                                             |    🟡    |  S   |
| M10       | godoc across the public API                                  |    🟡    |  M   |
| M11       | CI/CD + linter + coverage                                    |    🟡    |  M   |
| M12       | Semver + CHANGELOG                                           |    🟢    |  S   |
| M13       | API expansion (Equal / Compare / arithmetic)                 |    🟢    |  M   |
| M14       | `dto` evolution (Result, PaginatedResult)                    |    🟢    |  M   |

Sizes: S ≈ up to one day, M ≈ two–three days, L ≈ up to a week.

**Minimum work to reach the public release:** M1..M3 + M9 + M11, and ideally
M10. M1..M3 (and M5, M7) are now done — the remaining release-blocking work
is M9 (README overhaul) and M11 (CI pipeline).