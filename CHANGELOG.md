# Changelog

All notable changes to `gxsrvs/dtx` are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] — 2026-04-29

### Added

- **Comparison family** on every date/time type: `Before(other)`,
  `After(other)`, `Equal(other)`, `Compare(other) int`. The argument
  matches the receiver's type (not `time.Time`). Nullable variants
  use **sortable NULL semantics**: NULL is strictly less than any
  valid value, two NULLs compare equal, `Compare` returns
  `-1 / 0 / +1`. See `README.md` § «Methods on date / time types»
  for the full table.

- **Arithmetic** on every date/time type where applicable:

  - `Add(d time.Duration) X` — shift by duration. NULL propagates on
    nullable variants (NULL receiver returns NULL). On time-only
    types (`OffsetTime`, `LocalTime`) the result may rotate past
    midnight; the date is not part of the serialised form, so consume
    with care.
  - `AddDate(years, months, days int) X` — calendar shift. Not
    provided on time-only types.
  - `Sub(other X) time.Duration` on **not-null** types.
  - `SubOk(other NullX) (time.Duration, bool)` on **null** types —
    returns `(0, false)` if either operand is NULL, distinguishing
    NULL from a `0`-difference.
  - `Truncate(d) X` on every datetime/time type. Not on `Date`
    (already day-precision).

- **Unix family**:

  - Not-null types: `Unix() int64`, `UnixMilli()`, `UnixMicro()`,
    `UnixNano()` on `Date`, `OffsetDateTime`. (`Date.Unix*` reflects
    the underlying `time.Time`, which is conventionally midnight UTC
    for a Date value.)
  - Null types: `UnixOk() (int64, bool)`, `UnixMilliOk`, `UnixMicroOk`,
    `UnixNanoOk` — `bool` is `true` iff the receiver is valid.
    Replaces the legacy `NullOffsetDateTime.Unix() int64` (see
    `Removed` below).

- **Component getters** on not-null types: `Year()`, `Month()`,
  `Day()`, `YearDay()`, `Hour()`, `Minute()`, `Second()`,
  `Nanosecond()` where each component physically exists. Not added
  on null types — caller uses `.Val.X()` after `if .Valid`.
  `LocalDateTime` / `LocalTime` continue to expose components as
  struct fields directly (e.g. `ldt.Year`); only `YearDay()` is added
  there as a computed method.

- **`UTC()` shorthand** for `In(time.UTC)` on `OffsetDateTime` /
  `OffsetTime` and their nullable variants.

- **In/UTC absence on Local* and Date** — explicit godoc note added
  to the type-level documentation explaining why these types do not
  provide `In(loc)` or `UTC()` (they carry no timezone by design).

- `In(loc *time.Location)` on `OffsetDateTime`, `OffsetTime`,
  `NullOffsetDateTime` and `NullOffsetTime`. Mirrors `time.Time.In`:
  the absolute instant is preserved, only the displayed offset / wall
  clock changes. NULL propagates through the nullable variants.

- Alternative `*FromTime(time.Time)` constructors for every nullable
  date/time wrapper:
  - `NullDateFromTime`,
  - `NullOffsetDateTimeFromTime`,
  - `NullOffsetTimeFromTime`,
  - `NullLocalDateTimeFromTime`,
  - `NullLocalTimeFromTime`.

  Use them when you start from a `time.Time`. The primary
  `NewNullX(value X)` constructors continue to take the wrapper type.

### Changed

**BREAKING**: signature change on existing comparison methods. The
arguments now take the receiver's own type instead of `time.Time`:

| Method | Before | After |
| --- | --- | --- |
| `OffsetDateTime.Before/After` | `(time.Time) bool` | `(OffsetDateTime) bool` |
| `NullOffsetDateTime.Before/After` | `(time.Time) bool` | `(NullOffsetDateTime) bool` |

Migration: wrap the `time.Time` argument:

```go
// before
nodt.Before(time.Now())

// after
nodt.Before(NullOffsetDateTimeFromTime(time.Now()))
```

**BREAKING**: nullable date/time wrappers now hold the wrapper type
in their `Val` field, mirroring how `NullLocalDateTime` /
`NullLocalTime` already worked:

| Type                  | Before              | After              |
| --------------------- | ------------------- | ------------------ |
| `NullDate.Val`        | `time.Time`         | `Date`             |
| `NullOffsetDateTime.Val` | `time.Time`       | `OffsetDateTime`   |
| `NullOffsetTime.Val`  | `time.Time`         | `OffsetTime`       |

Wire format (Marshal / Scan / Value) is unchanged — only the in-memory
field type changes. Migration for callers reading `.Val` directly:

- `n.Val = t.Now()` → `n = NullDateFromTime(t.Now())` (or
  `n.Val = Date(t.Now())`).
- `n.Val.Year() / .Hour() / .Equal(t) / .Format(layout)` and friends
  → `n.Val.AsTime().Year()` etc., since `Date` / `OffsetDateTime` /
  `OffsetTime` don't inherit `time.Time` methods. Curated helpers
  (`ToString`, `In`, `Before`, `After`, `Unix` family) are exposed on
  the wrapper.

### Fixed

- `types.IsEmpty` now recognises `sql.NullByte` (added to the standard
  library in Go 1.17). Previously a `sql.NullByte{Valid: false}`
  argument fell through the type switch and was reported as non-empty.

### Removed

- `NullOffsetDateTime.Unix() int64` (introduced in v0.1.1) — its
  return-`0`-for-NULL semantics collided with a valid epoch instant.
  Replaced by `UnixOk() (int64, bool)`. Migration: switch to
  `UnixOk()`, or use `nv.Val.AsTime().Unix()` if NULL is guaranteed
  not to occur.

## [0.2.1] — 2026-04-28

### Fixed

- `OffsetDateTime`, `OffsetTime` and their nullable variants now parse
  the RFC 3339 `Z` suffix on `UnmarshalJSON` / `*FromString`. Before
  this fix only `±HH:MM` and `±HHMM` designators were recognised, so
  any input like `"2024-07-11T04:37:00Z"` or
  `"2024-07-11T04:37:00.123Z"` failed with
  `cannot parse time from …Z`. The library has always **emitted** the
  `Z` form for UTC values via `time.Format("…Z07:00")`, which broke
  the `Marshal → Unmarshal` round-trip for any UTC datetime. The fix
  is in `ParseTimezoneExtended`: a trailing `Z` is now stripped and
  treated as `time.UTC`.

## [0.2.0] — 2026-04-22

### Changed

**BREAKING**: public-API renames to follow the Go acronym convention
(capitalise initialisms — `revive`'s `var-naming` rule). Behaviour is
unchanged; this release is a pure rename.

| Before                                | After                                 |
| ------------------------------------- | ------------------------------------- |
| `types.NullUuid`                      | `types.NullUUID`                      |
| `types.NewNullUuid`                   | `types.NewNullUUID`                   |
| `types.NewNullUuidEmpty`              | `types.NewNullUUIDEmpty`              |
| `types.NullUuidFromString`            | `types.NullUUIDFromString`            |
| `utils.ToJson`                        | `utils.ToJSON`                        |
| `utils.LoadObjectFromJson`            | `utils.LoadObjectFromJSON`            |
| `utils.LoadCollectionFromJson`        | `utils.LoadCollectionFromJSON`        |
| `utils.LoadObjectFromJsonFile`        | `utils.LoadObjectFromJSONFile`        |
| `utils.LoadCollectionFromJsonFile`    | `utils.LoadCollectionFromJSONFile`    |

A blanket `sed -i 's/NullUuid/NullUUID/g; s/FromJson/FromJSON/g; s/ToJson/ToJSON/g'`
over the consumer codebase will migrate callers.

### Fixed

- `NSFromString` collapsed a redundant `valid := true` / `if s == ""`
  pattern into `Valid: s != ""` (`staticcheck QF1007`).
- `TestAssembleDateTimeTZ_InvalidZone` discarded its first return
  through a dead assignment; switched to `_` (`staticcheck SA4006`).
- CI: `golangci/golangci-lint-action` bumped `v6 → v7` with
  `install-mode: goinstall`, so the linter is rebuilt from source
  with the project's Go toolchain and stays compatible with
  `.golangci.yml` (v2 schema).

## [0.1.1] — 2026-04-22

### Added

- `Unix()` on `OffsetDateTime` and `NullOffsetDateTime`. Returns the
  underlying instant as a Unix timestamp (seconds since 1970-01-01
  UTC); the nullable variant returns `0` for NULL.

### Fixed

- `NullDecimal` and `NullUuid` now implement
  `database/sql/driver.Valuer`. Previously these two types satisfied
  only `sql.Scanner`, which caused drivers that rely on `Valuer` to
  serialise query parameters (notably `pgx` in simple-query mode) to
  fail when a `NullDecimal` or `NullUuid` was passed in. NULL emits
  `(nil, nil)`; valid values delegate to `decimal.Decimal.Value` and
  `uuid.UUID.Value` respectively.

## [0.1.0] — 2026-04-19

Initial public release.

[Unreleased]: https://github.com/gxsrvs/dtx/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/gxsrvs/dtx/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/gxsrvs/dtx/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/gxsrvs/dtx/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/gxsrvs/dtx/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/gxsrvs/dtx/releases/tag/v0.1.0
