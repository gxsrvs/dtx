# gxsrvs/dtx types library

[![Go Reference](https://pkg.go.dev/badge/github.com/gxsrvs/dtx.svg)](https://pkg.go.dev/github.com/gxsrvs/dtx)
[![Go Report Card](https://goreportcard.com/badge/github.com/gxsrvs/dtx)](https://goreportcard.com/report/github.com/gxsrvs/dtx)
[![CI](https://github.com/gxsrvs/dtx/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/gxsrvs/dtx/actions/workflows/ci.yml)
[![Coverage](https://codecov.io/gh/gxsrvs/dtx/branch/master/graph/badge.svg)](https://codecov.io/gh/gxsrvs/dtx)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A Go library providing nullable type wrappers for database operations with JSON marshaling/unmarshaling support.

## Installation

```bash
go get github.com/gxsrvs/dtx
```

## Features

This library provides not-null and nullable variants of common Go types
that are safe for database operations and JSON serialization. Each
not-null type shares its serializer with the nullable counterpart, so
required and optional fields render in identical formats.

More importantly, the library sharpens the data model itself: every
optional field carries a single explicit `Valid` flag instead of the
three-way ambiguity (`nil`, zero value, SQL `NULL`) that `*T` fields
expose by default. DTOs read unambiguously and serialise to the same
shape across Go, JSON, and SQL — one nullable concept, one spelling,
one contract.

### Date and time

| Not-null            | Nullable                | TZ?  | JSON format                          |
| ------------------- | ----------------------- | :--: | ------------------------------------ |
| `Date`              | `NullDate`              |  —   | `2006-01-02`                         |
| `OffsetTime`        | `NullOffsetTime`        | yes  | `15:04:05[.fff…]Z` / `±HH:MM`        |
| `OffsetDateTime`    | `NullOffsetDateTime`    | yes  | RFC 3339 (`2006-01-02T15:04:05…Z`)   |
| `LocalTime`         | `NullLocalTime`         |  no  | `15:04:05[.fff…]`                    |
| `LocalDateTime`     | `NullLocalDateTime`     |  no  | `2006-01-02T15:04:05[.fff…]`         |

`Local*` types reject inputs that carry a timezone designator. `Offset*`
types round-trip the offset (`Z` for UTC, `+HH:MM` / `-HH:MM` otherwise).

#### Accepted input / emitted output strings <sup>\*</sup>

The not-null and nullable variants share a single (de)serialiser, so the
input grammar and output shape are identical for each pair.

| Type pair | Accepted on input (Unmarshal / Parse / Scan-from-string)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   | Emitted on output (Marshal / ToString) |
| --- |--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------| --- |
| `Date` / `NullDate` | `2006-01-02` (ISO 8601), `02.01.2006` (dd.MM.yyyy).<br>Examples:<ul><li>`1961-04-12` — ISO 8601</li><li>`18.03.1965` — dd.MM.yyyy</li></ul>                                                                                                                                                                                                                                                                                                                                                                                                                                                | `2006-01-02`.<br>Example: `1961-04-12` |
| `OffsetTime` / `NullOffsetTime` | `15:04`, `15:04:05`, `15:04:05.fff`, `15:04:05.ffffff`, `15:04:05.fffffffff`, each optionally followed by a TZ designator (`Z`, `±HH:MM`, or `±HHMM`); a naked time is parsed in `time.Local`.<br>Examples:<ul><li>`06:07` — HH:MM, no TZ</li><li>`06:07:00` — no TZ</li><li>`02:56:15.123Z` — ms + UTC</li><li>`09:07:00+03:00` — `±HH:MM`</li><li>`11:34:51+0300` — `±HHMM`</li><li>`02:56:15.123456789-05:00` — ns + non-UTC</li></ul>                                                                                                                                                  | `15:04:05[.fff…]Z` for UTC, `15:04:05[.fff…]±HH:MM` otherwise (trailing zeros trimmed).<br>Examples:<ul><li>`20:17:40Z` — ns=0, UTC</li><li>`09:07:00+03:00` — ns=0</li><li>`02:56:15.123Z` — ms</li><li>`02:56:15.123456789Z` — full ns</li></ul> |
| `LocalTime` / `NullLocalTime` | `15:04`, `15:04:05`, `15:04:05.fff…` — **no** TZ designator allowed.<br>Examples:<ul><li>`06:07` — HH:MM</li><li>`06:07:00`</li><li>`02:56:15.123` — ms</li><li>`02:56:15.123456789` — full ns</li></ul>                                                                                                                                                                                                                                                                                                                                                                                   | `15:04:05[.fff…]` (trailing zeros trimmed).<br>Examples:<ul><li>`06:07:00` — ns=0</li><li>`02:56:15.123` — ms</li><li>`02:56:15.123456789` — full ns</li></ul> |
| `OffsetDateTime` / `NullOffsetDateTime` | Date (`2006-01-02` or `02.01.2006`) + `T` or space + time (`15:04`, `15:04:05`, `15:04:05.fff…`), with optional TZ designator (`±HH:MM` or `±HHMM`; an optional single space may precede it); a naked datetime is parsed in `time.Local`.<br>Examples:<ul><li>`1961-04-12T06:07:00Z` — `T` + UTC</li><li>`1969-07-20 20:17:40+0000` — space + `±HHMM`</li><li>`1965-03-18 11:34:51.000 +0300` — space before TZ</li><li>`18.03.1965 11:34:51 +03:00` — dd.MM.yyyy</li><li>`1969-07-21T02:56:15.123456+00:00` — μs</li><li>`1961-04-12T09:07:00` — no TZ (parsed in `time.Local`)</li></ul> | `2006-01-02T15:04:05[.fff…]Z` for UTC, `2006-01-02T15:04:05[.fff…]±HH:MM` otherwise.<br>Examples:<ul><li>`1969-07-20T20:17:40Z` — ns=0, UTC</li><li>`1961-04-12T09:07:00+03:00` — ns=0</li><li>`1969-07-21T02:56:15.123Z` — ms</li><li>`1969-07-21T02:56:15.123456789Z` — full ns</li></ul> |
| `LocalDateTime` / `NullLocalDateTime` | `2006-01-02T15:04:05[.fff…]` or `2006-01-02 15:04:05[.fff…]` — **no** TZ designator allowed.<br>Examples:<ul><li>`1961-04-12T06:07:00`</li><li>`1961-04-12 06:07:00` — space separator</li><li>`1969-07-21T02:56:15.123` — ms</li><li>`1969-07-21T02:56:15.123456789` — full ns</li></ul>                                                                                                                                                                                                                                                                                                  | `2006-01-02T15:04:05[.fff…]` (trailing zeros trimmed).<br>Examples:<ul><li>`1961-04-12T06:07:00` — ns=0</li><li>`1969-07-21T02:56:15.123` — ms</li><li>`1969-07-21T02:56:15.123456789` — full ns</li></ul> |

Notes:

- The `[.fff…]` notation in the output column means **the fractional part
  is optional and variable-length**: a single call to `ToString` /
  `MarshalJSON` picks exactly one shape per value, driven by the
  nanoseconds stored in the underlying `time.Time`. The canonical layout
  uses Go's `.999999999` verb, which omits the dot entirely when the
  fraction is zero and otherwise trims trailing zeros — so an `ms`-precise
  value emits three digits, a `μs`-precise value emits six, a `ns`-precise
  value emits nine, and a value with zero nanoseconds emits none.
- Fractional seconds are optional on input and preserved up to nanosecond
  precision; trailing zeros are trimmed on output.
- A `Local*` parser fails fast on any TZ designator (`Z`, `z`, `+HH:MM`,
  `-HH:MM`, `+HHMM`, `-HHMM`). Use the `Offset*` variant when the source
  carries an offset.
- Nullable variants additionally accept the JSON token `null` on
  `UnmarshalJSON` and SQL `NULL` on `Scan`; their not-null counterparts
  reject both.

<sup>\*</sup> **A note on the dates used in the examples above.** Wherever
a concrete calendar date appears in this section — as a parser input, a
canonical output, or a column illustration — the library uses milestones
of crewed and uncrewed spaceflight instead of arbitrary numbers. The
choice is purely decorative: parsing and formatting behaviour is
independent of the underlying values, and time-of-day components are
occasionally rounded for readability. The milestones referenced:

- `1957-10-04` — *Sputnik 1*, the first artificial Earth satellite.
- `1961-04-12` — Yuri Gagarin aboard *Vostok 1*, the first human
  spaceflight (lift-off at 06:07 UTC / 09:07 Moscow time from Baikonur).
- `1965-03-18` — Alexei Leonov performs the first extravehicular activity
  (EVA, i.e. spacewalk) from *Voskhod 2* (11:34 Moscow time).
- `1969-07-20` — Apollo 11 lunar module *Eagle* touches down in the Sea of
  Tranquillity at 20:17 UTC.
- `1969-07-21` — Neil Armstrong steps onto the lunar surface at 02:56 UTC
  (the instant fell late on July 20 in U.S. time zones — a useful
  reminder that a date depends on the zone you read it in).
- `1975-07-17` — the Apollo–Soyuz Test Project achieves the first
  international crewed docking in orbit.
- `2021-04-19` — *Ingenuity* performs the first powered, controlled flight
  of an aircraft on another planet (Mars).

### Other nullable wrappers

Beyond the date/time family, the library ships nullable wrappers over
Go primitives and a couple of value types from third-party packages
(`google/uuid`, `shopspring/decimal`). All wrappers share the same
conventions:

- a `Val` field holding the underlying value, plus a `Valid` flag;
- `Scan` delegates to the matching `sql.Null*` (or `sql.NullString`
  for text-encoded values), so a SQL `NULL` reliably yields
  `Valid == false`;
- `Value` returns `(nil, nil)` when invalid;
- `MarshalJSON` emits the JSON token `null` when invalid;
- `UnmarshalJSON` accepts both `null` and an empty input as NULL;
- `ToString` returns `""` when invalid;
- `IsEmpty()` is shorthand for `!Valid`;
- `New<T>(v)` constructs a valid wrapper, `New<T>Empty()` an invalid
  one;
- `<T>FromString(*string)` (where applicable) parses the underlying
  type from a `*string`, treating `nil`, `""`, and the case-
  insensitive tokens `"null"` / `"nil"` as NULL — a parse error also
  collapses to NULL rather than bubbling up.

| Type          | Underlying Go value                       | Valid JSON shape                            | Valid `ToString`                              | `Scan` source                       | Type-specific notes |
| ------------- | ----------------------------------------- | ------------------------------------------- | --------------------------------------------- | ----------------------------------- | ------------------- |
| `NullString`  | `string`                                  | JSON string                                 | the underlying string                         | `sql.NullString`                    | A valid empty string `""` is preserved as **valid** (`Valid == true`). `ToString` cannot tell a valid `""` from NULL — read the `Valid` field to disambiguate. The helper `NSFromString(s)` returns an `sql.NullString` that treats `""` as NULL, when you need that flavour at the `database/sql` boundary. |
| `NullBool`    | `bool`                                    | `true` / `false`                            | `"true"` / `"false"` (`fmt %t`)               | `sql.NullBool`                      | `UnmarshalJSON` accepts both the JSON literals `true` / `false` and the same tokens wrapped in quotes (`"true"` / `"false"`). |
| `NullInt16`   | `int16`                                   | JSON number                                 | decimal form (`fmt %d`)                       | `sql.NullInt16`                     | `Value` widens to `int64` per the `database/sql` contract. `UnmarshalJSON` accepts JSON numbers and number-strings. |
| `NullInt32`   | `int32`                                   | JSON number                                 | decimal form (`fmt %d`)                       | `sql.NullInt32`                     | Same widening + parsing rules as `NullInt16`. |
| `NullInt64`   | `int64`                                   | JSON number                                 | decimal form (`fmt %d`)                       | `sql.NullInt64`                     | Same parsing rules as `NullInt16`. |
| `NullFloat`   | `float64`                                 | JSON number                                 | `fmt %f` (six fractional digits)              | `sql.NullFloat64`                   | If exact decimal serialisation matters (money, ledger entries), prefer `NullDecimal` — `float64` is binary and cannot represent `0.1` exactly. |
| `NullDecimal` | `decimal.Decimal` (`shopspring/decimal`)  | JSON number (the form `decimal.Decimal` emits) | `decimal.Decimal.String()` — no trailing zeros | `sql.NullString` (text from the driver) | The decimal crosses the driver boundary as text to preserve precision. The helper `MulNullDecimals(a, b)` multiplies two `NullDecimal`s and returns NULL if either operand is NULL. |
| `NullUUID`    | `uuid.UUID` (`google/uuid`)               | JSON string in canonical `8-4-4-4-12` hex form | canonical `8-4-4-4-12` hex form               | `sql.NullString` (text from the driver) | Parsed via `uuid.Parse`, which accepts the canonical form, the bracketed `{…}` form, and the URN `urn:uuid:…` form. |

### Constructors for date / time types

Every nullable wrapper has a primary constructor that takes its
underlying type, and an alternative `*FromTime` constructor for
callers that start from a `time.Time`. The not-null types follow the
same shape: `Date` / `OffsetDateTime` / `OffsetTime` accept a
`time.Time` directly, while `LocalDateTime` / `LocalTime` are
field-based with a separate `*FromTime` helper (the wall-clock
representation has no zone, so the time.Time conversion is an
explicit step).

| Type                 | Primary                                           | From `time.Time`                            | From `string`                                       |
| -------------------- | ------------------------------------------------- | ------------------------------------------- | ------------------------------------------------- |
| `Date`               | `NewDate(time.Time)`                              | — primary already takes `time.Time`         | `DateFromString(string) (Date, error)`            |
| `OffsetDateTime`     | `NewOffsetDateTime(time.Time)`                    | — primary already takes `time.Time`         | `OffsetDateTimeFromString(string) (…, error)`     |
| `OffsetTime`         | `NewOffsetTime(time.Time)`                        | — primary already takes `time.Time`         | `OffsetTimeFromString(string) (…, error)`         |
| `LocalDateTime`      | `NewLocalDateTime(year, month, day, h, m, s, ns)` | `LocalDateTimeFromTime(time.Time)`          | `LocalDateTimeFromString(string) (…, error)`      |
| `LocalTime`          | `NewLocalTime(hour, minute, second, ns)`          | `LocalTimeFromTime(time.Time)`              | `LocalTimeFromString(string) (…, error)`          |
| `NullDate`           | `NewNullDate(Date)`                               | `NullDateFromTime(time.Time)`               | `NullDateFromString(*string)`                     |
| `NullOffsetDateTime` | `NewNullOffsetDateTime(OffsetDateTime)`           | `NullOffsetDateTimeFromTime(time.Time)`     | `NullOffsetDateTimeFromString(*string)`           |
| `NullOffsetTime`     | `NewNullOffsetTime(OffsetTime)`                   | `NullOffsetTimeFromTime(time.Time)`         | `NullOffsetTimeFromString(*string)`               |
| `NullLocalDateTime`  | `NewNullLocalDateTime(LocalDateTime)`             | `NullLocalDateTimeFromTime(time.Time)`      | `NullLocalDateTimeFromString(*string)`            |
| `NullLocalTime`      | `NewNullLocalTime(LocalTime)`                     | `NullLocalTimeFromTime(time.Time)`          | `NullLocalTimeFromString(*string)`                |

Note the asymmetry in the **From string** column: not-null parsers
take `string` and return `(T, error)` — a malformed input is a hard
failure. Nullable parsers take `*string` and return `NullT` with no
`error` — `nil` pointer, empty string, the tokens `null` / `nil`
(case-insensitive), or any parse error all collapse to an invalid
`NullT`. See the type's godoc for the accepted input formats.

Every nullable wrapper additionally has `NewNullXEmpty()` for the
explicit NULL case.

### Methods on date / time types

Every date/time type implements a curated subset of `time.Time`'s
behavior. ✓ — provided, ✗ — not provided (with a reason), ⚠ — provided
with a caveat documented in godoc.

| Method                         | Offset*  | OffsetTime | LocalDateTime | LocalTime | Date     |
| ------------------------------ | :------: | :--------: | :-----------: | :-------: | :------: |
| `In(loc)` / `UTC()`            | ✓        | ✓          | ✗ no zone     | ✗ no zone | ✗ no time |
| `Add(d)`                       | ✓        | ⚠ wraps    | ✓             | ⚠ wraps   | ✓        |
| `AddDate(y,m,d)`               | ✓        | ✗ no date  | ✓             | ✗ no date | ✓        |
| `Sub(other)` / `SubOk(other)`  | ✓        | ✓          | ✓             | ✓         | ✓        |
| `Truncate(d)`                  | ✓        | ✓          | ✓             | ✓         | ✗ day-precision |
| `Unix*()` / `Unix*Ok()`        | ✓        | ✗ no date  | ✗ needs zone  | ✗         | ✓ midnight UTC |
| `Year/Month/Day/YearDay`       | ✓        | ✗          | fields + `YearDay()` | ✗  | ✓        |
| `Hour/Minute/Second/Nanosecond` | ✓       | ✓          | fields        | fields    | ✗        |
| `Before/After/Equal/Compare`   | ✓        | ✓          | ✓             | ✓         | ✓        |

Notes:

- **Sortable NULL semantics** for the comparison family on nullable
  types: NULL is strictly less than any valid value, two NULLs are
  equal. `Compare` returns `-1 / 0 / +1` accordingly.
- **`Sub` is split** — not-null types expose `Sub(other) Duration`;
  nullable types expose `SubOk(other) (Duration, bool)` so a NULL
  operand never collides with a `0`-difference.
- **`Unix*` family** — same split: not-null returns `int64`, nullable
  returns `(int64, bool)` via `UnixOk` / `UnixMilliOk` / `UnixMicroOk`
  / `UnixNanoOk`.
- **`In/UTC` not on Local* / Date** — those types carry no timezone
  by design; zone reinterpretation is meaningless. Materialise via
  `LocalDateTime.ToTime(loc)` / `LocalTime.ToTime()` if needed.
- **`Add` on time-only types** (`OffsetTime`, `LocalTime`) — does not
  enforce modulo-24h. A duration past midnight rotates the implicit
  date in the underlying `time.Time`; the date is not part of the
  serialised form, so consume the result with that in mind.
- **Component fields on Local***  — `LocalDateTime` and `LocalTime`
  expose `Year/Month/Day/Hour/Minute/Second/Nanosec` as struct fields
  directly (no method call). Only `LocalDateTime.YearDay()` is a
  method (it is computed, not stored).

## Usage

### Basic Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/gxsrvs/dtx/types"
)

func main() {
    // Create a nullable string with value
    ns := types.NewNullString("Hello, World!")
    fmt.Println("Value:", ns.ToString())
    fmt.Println("Is Empty:", ns.IsEmpty())
    
    // Create an empty nullable string
    emptyNs := types.NewNullStringEmpty()
    fmt.Println("Empty Value:", emptyNs.ToString())
    fmt.Println("Is Empty:", emptyNs.IsEmpty())
    
    // JSON marshaling
    jsonData, _ := json.Marshal(ns)
    fmt.Println("JSON:", string(jsonData))
    
    // JSON unmarshaling
    var unmarshaled types.NullString
    json.Unmarshal([]byte(`"Test Value"`), &unmarshaled)
    fmt.Println("Unmarshaled:", unmarshaled.ToString())
}
```

### Date Handling

```go
package main

import (
    "fmt"
    "time"
    "github.com/gxsrvs/dtx/types"
)

func main() {
    // Create nullable date from a time.Time
    now := time.Now()
    nd := types.NullDateFromTime(now)
    fmt.Println("Date:", nd.ToString())

    // Parse date from string (supports ISO and Russian formats)
    dateStr := "1961-04-12" // Vostok 1 — first human spaceflight
    parsed, err := types.ParseDateFromString(dateStr)
    if err == nil {
        nd2 := types.NullDateFromTime(*parsed)
        fmt.Println("Parsed Date:", nd2.ToString())
    }
}
```

### Database Integration

All types implement the `sql/driver.Valuer` and `sql.Scanner` interfaces for seamless database integration:

```go
package main

import (
    "database/sql"
    "github.com/gxsrvs/dtx/types"
)

type User struct {
    ID       int64             `db:"id"`
    Name     types.NullString  `db:"name"`
    Birthday types.NullDate    `db:"birthday"`
    Active   types.NullBool    `db:"active"`
}

func queryUser(db *sql.DB, id int64) (*User, error) {
    user := &User{}
    err := db.QueryRow("SELECT id, name, birthday, active FROM users WHERE id = ?", id).
        Scan(&user.ID, &user.Name, &user.Birthday, &user.Active)
    return user, err
}
```

### Timezone conversion

`OffsetDateTime`, `OffsetTime` and their nullable variants expose
`In(loc *time.Location)`, mirroring `time.Time.In`. The absolute
instant is preserved; only the displayed offset and wall clock change.
NULL propagates through the nullable variants.

```go
import (
    "time"

    "github.com/gxsrvs/dtx/types"
)

odt, _ := types.OffsetDateTimeFromString("2024-07-11T10:00:00+04:00")
odt.In(time.UTC).ToString()
// → "2024-07-11T06:00:00Z"

utc, _ := types.OffsetDateTimeFromString("2024-07-11T06:00:00Z")
utc.In(time.FixedZone("UTC+04:00", 4*3600)).ToString()
// → "2024-07-11T10:00:00+04:00"
```

## Supported date formats

`Date` and `NullDate` accept the following input formats:

- ISO 8601: `2006-01-02`
- dd.MM.yyyy locale form: `02.01.2006`

Output is always ISO 8601.

`Offset*` types accept the corresponding date / time / datetime formats
with an optional timezone designator (`Z`, `+HH:MM`, or `+HHMM`). Output
always carries an explicit designator (`Z` for UTC, `+HH:MM` / `-HH:MM`
otherwise).

`Local*` types accept the same shapes but reject any timezone designator
on input and never emit one on output.

## Performance

The `Null*` wrappers carry a small overhead compared to plain
`*T` pointer fields: an extra `bool` per field, more allocations on the
JSON marshal path, and a denser struct layout. The `bench/` package
quantifies this end-to-end on a representative `Client` DTO (six
nullable fields across two structs) and compares the two styles head
to head:

| Op / profile         |  Null (ns) |  Ptr (ns) |  Null (B/op) | Ptr (B/op) | Null (allocs/op) | Ptr (allocs/op) |
| -------------------- | ---------: | --------: | -----------: | ---------: | ---------------: | --------------: |
| Marshal / AllValid   |       3593 |      2895 |          592 |        480 |               14 |               5 |
| Marshal / Mixed      |       1005 |       641 |          248 |        160 |                4 |               2 |
| Marshal / AllNull    |        377 |       302 |          128 |         64 |                2 |               2 |
| Unmarshal / AllValid |       8028 |      5300 |         1072 |        584 |               20 |              20 |
| Unmarshal / Mixed    |       3187 |      2761 |          632 |        424 |               13 |              12 |
| Unmarshal / AllNull  |        705 |       678 |          328 |        264 |                5 |               5 |

(Linux / Ryzen 7 8845H, three runs each; reproduce with
`go test ./bench/ -bench=BenchmarkClient -benchmem -run=^$`.)

The gap is real but modest — well under a microsecond per field on
typical DTOs. In return, defining the DTO becomes much more pleasant:

- **One nil-state instead of two.** A `*string` can be `nil` *or* point
  to `""`; a `NullString` has a single `Valid` flag, so call sites do
  not need to guess which sentinel a producer chose.
- **SQL `Scan` works out of the box.** A bare `*time.Time` does not
  implement `sql.Scanner`, so you would have to write a wrapper anyway —
  the `Null*` types delegate to the matching `sql.Null*` and honour
  the driver's NULL signalling reliably.
- **Compact, predictable JSON.** `NullDate` writes `"1961-04-12"` rather
  than RFC 3339 `"1961-04-12T00:00:00Z"`; the `Local*` and `Offset*`
  pairs each pick one canonical shape. On the `AllValid` profile this
  saves 30 bytes per record (248 B vs 278 B for `*time.Time`).
- **Symmetric `Marshal` / `Unmarshal` / `Scan` / `Value` semantics.**
  `null`, `""`, `"null"`, and SQL `NULL` collapse to the same
  invalid-wrapper state across every type, so DTO-level handling does
  not need special cases per field.

If the per-record cost matters in a hot path, mix and match: the
library does not force a choice. Pointer fields and `Null*` fields
coexist in the same struct, so leave the rare bottleneck as `*T` and
keep `Null*` for the rest of the DTO.

## Dependencies

- `github.com/google/uuid` — UUID handling for `NullUUID`.
- `github.com/shopspring/decimal` — arbitrary-precision decimals for `NullDecimal`.

Both are permissive-licensed (BSD-3-Clause and MIT respectively). See
[`THIRD_PARTY_LICENSES.md`](THIRD_PARTY_LICENSES.md) for the attribution
notices that consumers of this library must preserve when redistributing.

## Contributing

Contributions are welcome — please open an issue or a pull request.

**Project language: English.** All repository artefacts (source code, godoc
comments, Markdown documentation, commit messages, issue and pull request
descriptions, configuration files) must be written in English. This keeps the
project accessible to any Go developer who picks it up.

## License

This library is available under the MIT License.