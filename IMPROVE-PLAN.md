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

## M4. ✅ Consistent `ToString()`

Rules:

- `Valid == false` → `""` across every type,
- for date/time types — the **same** format as `MarshalJSON`.

Done:

- date/time family shares its formatters with `MarshalJSON`
  (e.g. `NullOffsetTime` routes through `formatOffsetTime`),
- `NullFloat.ToString()` returns `""` for invalid values, aligning with
  every other nullable wrapper (covered by `null_float_test.go`).

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

## M8. ✅ `utils.ToJson` returns an error

Done. The signature now is:

```go
func ToJson(entity interface{}) (string, error)
```

`json.Marshal` failures are wrapped with `fmt.Errorf("failed to marshal
interface to json: %w", err)` and returned to the caller instead of being
logged and swallowed. `utils/json_test.go` covers the error path (channel
marshalling) explicitly.

---

## M9. 🟡 Add `LICENSE` + overhaul `README.md`

`LICENSE` (MIT) and `THIRD_PARTY_LICENSES.md` are already in place. The
remaining work is the README overhaul per `.ai/REVIEW.md` §7.

Done:

- ✅ all five badges in place (`pkg.go.dev`, Go Report Card, CI,
  coverage, MIT license) — `README.md:3-7`. The MIT-license badge is
  a static shields.io image and renders immediately; CI / coverage
  badges will only render once M11 lands the workflow + Codecov
  upload; `pkg.go.dev` and Go Report Card light up after the first
  public push/tag,
- ✅ Installation section,
- ✅ Contributing section,
- ✅ English as the sole language of the docs (explicit note in the
  Contributing section),
- ✅ full documentation of every public type — the date/time pairs
  have a dedicated input/output grammar table; the simple-type family
  (`NullString`, `NullBool`, `NullInt{16,32,64}`, `NullFloat`,
  `NullDecimal`, `NullUuid`) now has a parallel table covering
  underlying Go value, JSON shape, `ToString` form, `Scan` source, and
  per-type idiosyncrasies, plus a paragraph stating the conventions
  shared across the family (`README.md` "Other nullable wrappers"
  section).

Still open:

- 🟡 dedicated sections Quickstart / SQL / JSON / TZ / Testing — at
  the moment Usage merges Quickstart + SQL + JSON examples, and there
  are no standalone JSON / TZ / Testing sections,
- 🟡 semver policy — `.ai/REVIEW.md` §6 already drafts it; needs to be
  distilled into a README section.

---

## M10. ✅ godoc coverage

Done. Every exported type, function, method, and constant in `types/`
and `utils/` now carries a godoc comment. Package-level overviews live
in `types/doc.go` and `utils/doc.go`. The nullable families
(`NullBool`, `NullInt{16,32,64}`, `NullFloat`, `NullString`,
`NullDecimal`, `NullUuid`, `NullDate`, and the four `Null*Local*` /
`Null*Offset*` datetime types) document their Scan/Value/Marshal
semantics and the `""` / `null` conventions in one voice; the not-null
`Date` / `LocalDateTime` / `LocalTime` / `OffsetDateTime` /
`OffsetTime` types document their shared serializer with the nullable
counterpart.

---

## M11. 🟡 CI/CD

Done:

- ✅ `.github/workflows/ci.yml` with two jobs: **lint** (golangci-lint
  v2, latest) and **test** (`go build ./...` + `go test -race
  -covermode=atomic -coverprofile=coverage.out ./...` + Codecov
  upload via `codecov/codecov-action@v5`),
- ✅ `.golangci.yml` (v2 schema) enabling `errcheck`, `govet`,
  `staticcheck`, `unused`, `revive` as linters and `gofmt`,
  `goimports` as formatters,
- ✅ `.gitignore` extended with `coverage.out` / `coverage.html`,
- ✅ Go version pinned to **1.26 only** — the minimum supported
  version, no matrix. We deliberately use 1.26-only features
  (`new(expr)` shorthand) and the `b.Loop()` benchmark idiom from
  1.24, so dropping the floor would force a refactor with no real
  upside for a fresh public library.

Still open:

- 🟡 Codecov badge will only render once the `gxsrvs/dtx` repo is
  added on `codecov.io` (a one-time web-UI step) and the first CI
  run uploads coverage. No token is needed for public repos.
- 🟡 Release workflow (goreleaser or `release-please`) — deferred to
  M12 alongside the broader versioning / `CHANGELOG.md` work, since
  the two are tightly coupled.

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

## Roadmap

| Milestone | Description                                                  | Priority | Size |
| --------- | ------------------------------------------------------------ | :------: | :--: |
| M1        | Non-null Date/OffsetTime/OffsetDateTime + shared serializer  |    ✅    |  M   |
| M2        | LocalDateTime/NullLocalDateTime + LocalTime/NullLocalTime    |    ✅    |  L   |
| M3        | `Scan` fix (use `sql.Null*.Valid`) across every nullable     |    ✅    |  M   |
| M4        | Consistent `ToString` — all nullable wrappers return ""      |    ✅    |  S   |
| M5        | Rename `Time*` / `DateTime*` → `Offset*` (breaking, pre-v1.0)|    ✅    |  S   |
| M6        | Receivers + remove `//goland:noinspection` — new files done  |    🟡    |  S   |
| M7        | Interface dispatch for `IsEmpty` / `ToString`                |    ✅    |  S   |
| M8        | `utils.ToJson` returns an error                              |    ✅    |  S   |
| M9        | LICENSE + README                                             |    🟡    |  S   |
| M10       | godoc across the public API                                  |    ✅    |  M   |
| M11       | CI/CD + linter + coverage                                    |    🟡    |  M   |
| M12       | Semver + CHANGELOG                                           |    🟢    |  S   |
| M13       | API expansion (Equal / Compare / arithmetic)                 |    🟢    |  M   |

Sizes: S ≈ up to one day, M ≈ two–three days, L ≈ up to a week.

**Minimum work to reach the public release:** M1..M3 + M9 + M11, and ideally
M10. M1..M3 (and M5, M7, M10) are now done. M11's CI pipeline (lint + test +
coverage) is wired up and stays 🟡 only because the Codecov web-UI sign-up and
the release workflow are still outstanding. Remaining release-blocking work is
the rest of M9 (Quickstart / SQL / JSON / TZ / Testing sections + semver
policy distillation).