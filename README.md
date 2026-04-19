# gxsrvs/dtx types library

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

- `NullString` — nullable string
- `NullBool` — nullable boolean
- `NullInt16`, `NullInt32`, `NullInt64` — nullable integers
- `NullFloat` — nullable float
- `NullDecimal` — nullable arbitrary-precision decimal (`shopspring/decimal`)
- `NullUuid` — nullable UUID (`google/uuid`)

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
    // Create nullable date
    now := time.Now()
    nd := types.NewNullDate(now)
    fmt.Println("Date:", nd.ToString())
    
    // Parse date from string (supports ISO and Russian formats)
    dateStr := "1961-04-12" // Vostok 1 — first human spaceflight
    parsed, err := types.ParseDateFromString(dateStr)
    if err == nil {
        nd2 := types.NewNullDate(*parsed)
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

## Dependencies

- `github.com/google/uuid` — UUID handling for `NullUuid`.
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