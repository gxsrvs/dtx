# `gxsrvs/dtx` — Project Review

This document captures the current state of the codebase before the library is
released as public open-source software. It lists the strengths and weaknesses
that shaped the follow-up work planned in `IMPROVE-PLAN.md` and the test
coverage roadmap in `TESTS-PLAN.md`.

## 1. Purpose

`dtx` is a collection of nullable wrappers around Go primitives, intended to
describe DTO structures that travel both through `database/sql` and through
JSON. Typical consumers include:

- HTTP/REST services: `DB → Go struct → JSON → client` (and back);
- message-driven services: `Kafka / NATS / RabbitMQ → Go struct → business logic
  → Kafka / NATS / RabbitMQ` (and back);
- any pipeline where the same Go struct must survive JSON serialisation and
  SQL scanning, and where `NULL` must be represented faithfully at every hop.

## 2. Repository layout

```
dtx/
├── dto/                 — generic DataPackage[T] container
├── types/               — nullable wrappers + base DateTime / TimeOnly types
├── utils/               — JSON-serialisation helpers
├── README.md            — user-facing documentation
└── go.mod               — Go 1.26.1; deps: google/uuid, shopspring/decimal
```

## 3. Strengths

### 3.1. Functional
- **Typical SQL-friendly type set is covered:** `NullString`, `NullBool`,
  `NullInt16/32/64`, `NullFloat`, `NullDecimal`, `NullDate`, `NullIsoDate`,
  `NullDateTime`, `NullTime`, `NullUuid`.
- **Consistent constructor shape:** `New<T>(val)`, `New<T>Empty()`,
  `<T>FromString(*string)` — predictable for callers.
- **Consistent JSON marshaling:** invalid → `null`, valid → value in the
  documented format. A shared `nullJson` constant lives in
  `types/null_json.go`.
- **Complete `sql.Scanner` / `driver.Valuer` implementation** for every
  nullable type, so they drop straight into `Query`, `QueryRow`, `Scan`.
- **Multiple date formats** (`ISO 8601` plus the `dd.MM.yyyy` locale format)
  and a relatively rich timezone parser in `time_only.go`
  (`ParseTimezoneExtended`, `parseTimezoneWithColon`, `parseTimezone`).
- **Helpers in `utils/json.go`** for loading structs or collections from a
  JSON string or file (`LoadObjectFromJson`, `LoadCollectionFromJsonFile`).
- **Generic `StdDataPackage[T]`** in `dto/` that standardises list-shaped API
  responses.
- **Unit tests exist** for most public entry points; `go test ./...` is green
  and the module builds cleanly.

### 3.2. Organisational
- License declared MIT (in the README).
- Sensible package split (`types`, `utils`, `dto`).
- Minimal external dependencies: `google/uuid`, `shopspring/decimal`.

## 4. Weaknesses and risks

### 4.1. Semantic bugs / surprises
- **`Scan` in `NullDate`, `NullDateTime`, `NullTime`, `NullString`, `NullBool`
  and the numeric types decides `Valid` via `reflect.TypeOf(value) == nil`.**
  The database driver may hand in a typed nil (for example `([]byte)(nil)`)
  or a `sql.NullTime{Valid:false}`. The correct source of truth is the
  `Valid` field of the embedded `sql.Null*`. The current logic silently flips
  such inputs to `Valid=true` with a zero value — a real `NULL` can get
  turned into a zero value without error.
- **`DateTime.Scan` branches on `if reflect.TypeOf(value) == nil { … } else
  { … }` with identical bodies in both arms** (`types/datetime.go:74-79`).
  It is either dead code or unfinished — in both cases it confuses the reader.
- **`NullFloat.ToString()` returns `"null"` when invalid**
  (`types/null_float.go:49`); every other nullable wrapper returns an empty
  string in the same state. Inconsistent API.
- **`NullTime.ToString()` formats via `thisVal.Val.String()`**
  (`types/null_time.go:48`), but `MarshalJSON` uses `Format(time.TimeOnly)`.
  Two different string representations for the same value.
- **`TimeOnly.Value()` returns `thisVal` typed as `TimeOnly`**
  (`types/time_only.go:146`) — the SQL driver does not know it is a
  `time.Time` alias. It needs an explicit `time.Time(thisVal)` cast (or a
  formatted string, depending on the target column type).
- **The name `NullTime` is misleading.** In the standard library
  `sql.NullTime` is a nullable `timestamp` (date + time). In this package
  `NullTime` marshals as a time-of-day (`15:04:05`). For a public library
  this is a trap — a caller familiar with the `database/sql` package will
  misuse it.
- **Date parsing uses `time.Local`** (`ParseDateFromString`). On a server
  this is a common source of drift when the process is restarted in a
  different TZ.

### 4.2. Architectural
- **No non-null Date/Time/DateTime types that share the library’s serializer.**
  Today, a DTO that consumes this library has to use bare `time.Time` for
  required datetime fields (which marshals as stdlib RFC3339 with timezone),
  while optional ones use the library’s own formats. A `Date`/`NullDate`
  pair (and `Time`/`NullTime`, `DateTime`/`NullDateTime`) built on one shared
  serializer would remove that asymmetry.
- **No `LocalDateTime` without timezone.** For business events that are
  naturally local (“contract signed on 2026-04-18 at 13:00” — a moment in
  the wall-clock time of an office, with no offset) a dedicated local
  datetime type is needed; `time.Time` always carries a zone and is the
  wrong modelling choice.
- **The non-null `DateTime` exists but the API is incomplete:** no
  `Value()`, no `IsEmpty`/`ToString` on the value, and `UnmarshalJSON`
  swallows `null` as a no-op (it mutates the receiver only on valid input).
- **`Emptiable` and `IsEmpty` in `types/emptiable.go`** use a long `type
  switch` that duplicates the list of types. Adding a new type needs edits
  in two places. An interface-based dispatch is simpler:
  `if e, ok := v.(Emptiable); ok { return e.IsEmpty() }`.
- **Same for `ToString` (`types/string_able.go`)** — another hand-rolled
  `type switch` over the full catalogue.
- **Mixed value- and pointer-receivers** on single types
  (`//goland:noinspection GoMixedReceiverTypes` appears dozens of times).
  Idiomatic Go prefers one receiver style per type.
- **`utils.ToJson` swallows the error** (`log.Println` + return `""`). That
  is not appropriate for a public library API.
- **`dto.DataObject = interface{}`** — a type alias for the empty
  interface gives neither type safety nor documentation value, and is
  unnecessary post-Go 1.18 when `any` and generics are available.

### 4.3. Code quality
- **Pervasive `//goland:noinspection GoMixedReceiverTypes`** — a symptom of
  IDE warnings being suppressed rather than fixed. A public library should
  have a clean listing.
- **Russian comments in a handful of files** (parser loops in
  `null_date.go`, `time_only.go`, etc.). For open-source work the public
  API, including comments, must be in English.
- **No `// Package ...` doc comment** in any package — godoc is going to
  look empty.
- **`NullString.Scan` differs from its siblings**: it round-trips through
  `json.Unmarshal`, while the other nullable types use `strings.Trim` and a
  direct assignment. That stylistic drift is worth removing in one direction.

### 4.4. Infrastructure
- ~~No `LICENSE` file~~ — added (`LICENSE`, MIT). Third-party attribution
  notices reproduced in `THIRD_PARTY_LICENSES.md`.
- **No CI** (GitHub Actions / lint / test). Nothing currently guarantees
  that `master` stays green.
- **No `CHANGELOG.md`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`.**
- **No release tags.** Commit messages mention `v0.1.0`, but `git tag` has
  nothing — consumers would have to rely on Go’s pseudo-versions.
- **`go 1.26.1` — the newest available.** For a public module the minimum
  supported Go version should be pinned lower (e.g. `go 1.22` or `go 1.23`)
  and documented.
- **`.gitignore` contains `/.claude/`, `/.ai/`, `CLAUDE.md`** — fine for a
  private workspace, but should be reviewed before public release.
- **Tests are mostly not table-driven** and coverage is uneven
  (`NullString`, `NullDecimal` are well covered; `NullBool`, `NullFloat`
  are not).

## 5. Public release readiness

| Criterion                                   | Status |
| ------------------------------------------- | :----: |
| `go build` / `go test ./...` succeed        |   ✅   |
| `LICENSE` file at repo root                 |   ✅   |
| Package-level godoc                         |   ❌   |
| Comments in English                         |   ⚠️   |
| Unit-test coverage ≥ 80%                    |   ⚠️   |
| CI / linter / coverage report               |   ❌   |
| `CHANGELOG.md`, `CONTRIBUTING.md`           |   ❌   |
| Semantic version tags in git                |   ❌   |
| Non-null Date/Time/DateTime + LocalDateTime |   ❌   |
| Consistent ToString / JSON representation   |   ⚠️   |

## 6. README suggestions

The file exists but has gaps:

1. Add badges (pkg.go.dev reference, Go Report Card, CI status, coverage,
   MIT license).
2. Document **every** exported type — the current README does not mention
   `NullUuid`, `NullDecimal`, `NullIsoDate`, `DateTime`, `TimeOnly`, or
   `StdDataPackage`.
3. Add a “Stability / Semver policy” section: what `v0.x` means and what
   freezes at `v1.0`.
4. Add short recipes:
   - a JSON round-trip,
   - `sql.Row.Scan` with nullable columns,
   - assembling a `DateTime` from `NullDate + NullTime + TZ` via
     `AssembleNullDateTimeTZ`,
   - using `StdDataPackage[T]` as a standard list-response envelope.
5. Clarify the **semantics of `NullTime`** (it is a time-of-day, not a
   `sql.NullTime` equivalent) — or rename the type, see `IMPROVE-PLAN.md`.
6. State the minimum Go version and support policy.
7. Add a “Contributing” section with a link to the issue tracker.
8. Translate every inline example and comment to English.