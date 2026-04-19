# `gxsrvs/dtx` — Testing Plan

The goal is to bring library test coverage up to public-release quality
(target ≥ 85% line coverage; JSON-null / SQL-NULL / invalid-input branches
are mandatory) and to formalise the coverage gate in CI.

Tests are plain `testing` + table-driven, without third-party frameworks.
Integration scenarios that require `database/sql` semantics use
`github.com/DATA-DOG/go-sqlmock`.

**Current status:** unit-test coverage is at **95.0%** of statements in
`types/`, **96.3%** in `utils/`. Phases T1, T3, T4, T5
and T6 are done; T2 (sqlmock integration), T7 (fuzzing), T8 (benchmarks)
and T9 (CI) are still open.

## 1. Principles

1. **Table-driven.** Each test is a `[]struct{ name string; in …; want …;
   wantErr bool }`.
2. **JSON round-trip.** `Marshal → Unmarshal → DeepEqual` for every nullable
   type.
3. **SQL round-trip.** `Value() → Scan(value)` against a fake row served by
   `sqlmock`.
4. **Edge cases are mandatory:**
   - `nil` input,
   - empty string,
   - `"null"` / `"nil"` (case-insensitive) where the API accepts it,
   - invalid format,
   - JSON literal `null` vs. JSON string `"null"` (they are different),
   - DB `NULL`,
   - numeric boundaries (int16 min/max, decimal overflow, large float).
5. **Time zones.** Every date/time test runs under at least two locales:
   `time.UTC` and `time.FixedZone("SampleTZ", 3*3600)`.
6. **Fuzzing.** Parsers for dates/times/timezones are exercised with
   `go test -fuzz`.
7. **Benchmarks** for hot paths (`MarshalJSON`, `UnmarshalJSON`,
   `ParseDateTimeFromString`).

## 2. Coverage, file by file

### 2.1. `types/null_string.go` — `null_string_test.go` ⚠️
Baseline plus `IsEmpty`, `ToString`, `Value`, `Scan`, JSON marshal/unmarshal
are covered. The Scan transition "valid → NULL" is covered in
`scan_reset_test.go`. Still open:
- UnmarshalJSON for escaped characters (`"a\"b"`, `"\u0026"`),
- Scan from `[]byte`, from `int` (error), from `time.Time` (error),
- round-trip for non-ASCII strings.

### 2.2. `types/null_bool.go` — `null_bool_test.go` ✅
Done. Covers constructors, `NullBoolFromString` (incl. `null`/`NIL`/garbage),
`IsEmpty`, `ToString`, `Value`, `Scan` (bool, int64, nil), `MarshalJSON`,
`UnmarshalJSON` (`true`/`false`/`null`/`"true"`/garbage).

### 2.3. `types/null_int16.go`, `null_int32.go`, `null_int64.go` ✅
Done. Each has constructors, `FromString` (incl. overflow for int16, `null`
sentinels), `FromNullString` (where present), `IsEmpty`, `ToString`,
`Value`, `Scan` (int64, string, nil), `MarshalJSON`, `UnmarshalJSON`
(numeric, stringified, `null`, garbage).

Still open: Scan from an explicitly unsupported type and `Scan` from
typed-nil drivers (`([]byte)(nil)`).

### 2.4. `types/null_float.go` — `null_float_test.go` ✅
Done. Covers constructors, `NullFloatFromString`, `IsEmpty`, `ToString`
(documents the M4 inconsistency where invalid → `"null"`), `Value`, `Scan`,
JSON marshal/unmarshal, and `math.MaxFloat64` boundary round-trip.

Still open: `+Inf` / `-Inf` / `NaN` — pending API decision whether they are
supported (see `IMPROVE-PLAN.md` M4).

### 2.5. `types/null_decimal.go` — `null_decimal_test.go` ⚠️
Done: constructors, `NullDecimalFromString`, `IsEmpty`, `Scan` (string /
`[]byte` / nil), JSON round-trip on `0.1 + 0.2`, marshal-empty → `null`,
unmarshal-null + unmarshal-garbage. `MulNullDecimals` already covered.

Still open:
- `Scan` from `float64` (error case),
- explicit assertion that `MarshalJSON` emits a bare number (no quotes) —
  for JS/TS compatibility.

### 2.6. `types/null_uuid.go` ✅
Done. Constructors, `NullUuidFromString` (incl. `null`/`NIL`/empty),
`IsEmpty`, `ToString`, `Scan` (string / `[]byte` / nil / invalid),
`MarshalJSON`, `UnmarshalJSON` (canonical, `null`, garbage).

Still open: braced (`{xxxxxxxx-…}`) and upper-case input variants.

### 2.7. `types/null_date.go` ✅
Done. Covers constructors, `FromString` (ISO + `dd.MM.yyyy`),
`IsEmpty`, `ToString`, `Value`, `Scan` (time.Time, nil),
`MarshalJSON`, `UnmarshalJSON` (incl. `null` and garbage).
`ParseDateFromString` is covered for both supported formats and the error
path.

Still open:
- explicit invalid-calendar inputs (`2023-13-01`, `31.02.2023`,
  `0000-00-00`, `2023/01/02`),
- `DateToString` across non-UTC time zones.

### 2.8. `types/{offset,null_offset}_datetime.go` ✅
Done (under the new naming `OffsetDateTime` / `NullOffsetDateTime`):
- `T`-separator parsing and JSON round-trip with TZ,
- UTC marshalling emits `Z`,
- not-null rejects JSON `null` and empty string,
- `Value()` / `Scan()` (incl. NULL-rejection on the not-null type and
  NULL → invalid on the nullable),
- `Before` / `After` on both types (nullable form returns false when
  invalid),
- `ToString`,
- `DateTimeToString` wrapper.

Still open:
- explicit space-separator parsing (`2023-05-20 15:30:45+03:00`),
- fractional-seconds round-trip (`2023-05-20T15:30:45.123Z`),
- Scan from `string` / `[]byte` (currently driver-side only).

### 2.9. `types/{offset,null_offset}_time.go` ✅
Done (under the new naming `OffsetTime` / `NullOffsetTime`):
- format round-trip with TZ (`+03:00`) and UTC (`Z`),
- `String()` / `ToString()` produce identical output,
- `Value`, `Scan` (time.Time, string, []byte, unsupported, invalid string),
- `MarshalJSON` / `UnmarshalJSON` (null rejected on not-null; null → invalid
  on nullable),
- `IsEmpty` on `NullOffsetTime`.
- `isDigitString`, `hasOffsetSuffix`, `ParseTimezoneExtended`
  (no-suffix / colon / no-colon forms) covered in `types_utils_test.go`.

Still open:
- direct tests for `parseTimezone` and `parseTimezoneWithColon` error paths
  (only happy paths are exercised through `ParseTimezoneExtended`),
- explicit `HH:MM:SS.ms` fractional-seconds round-trip.

### 2.10. `types/emptiable.go`, `types_utils.go`, `string_able.go` ✅
Done. `IsEmpty` and `ToString` are exercised on every library type (both
value and pointer kinds), incl. the `nil` and unknown-type fallbacks; the
`stdlib sql.Null*` branch is covered. `MaxDateTime` / `MinDateTime`,
`AssembleDateTime` (with and without explicit location) and
`AssembleNullDateTimeTZ` (all four valid/invalid combinations) are tested.

Still open:
- `AssembleDateTimeTZ` with a malformed-zone error (the parser is lenient
  and currently swallows most malformed inputs into `time.Local`).

### 2.11. `utils/json.go` ⚠️
Baseline exists; coverage at 96.3%. Still open:
- `LoadObjectFromJson` with a generic type that contains
  `NullString`/`NullDate` — confirm round-trip preservation,
- `LoadCollectionFromJsonFile` — empty file, non-array input (error),
- `ToJson` — error path is covered (channel marshalling); consider adding
  a nested-struct success case.

## 3. New types introduced by `IMPROVE-PLAN.md`

Each new type (`Date`, `Time`, `DateTime`, `LocalDateTime`,
`NullLocalDateTime`) gets a mandatory test set:
- constructors,
- JSON round-trip (including `null` for nullable forms),
- SQL round-trip (`Value() → Scan()`) through `sqlmock`,
- parsing every supported string format,
- parse errors on invalid input,
- timezone semantics (`LocalDateTime` must **not** attach a TZ on
  serialization),
- consistency with the nullable sibling (where applicable):
  `Date.ToString() == NullDate.ToString()` for the same underlying
  `time.Time`.

## 4. Integration tests

### 4.1. SQL round-trip (`types/integration_sql_test.go`)

In-process transition tests live in `types/scan_reset_test.go` and assert
that every nullable `Scan` resets cleanly from a previously-valid state when
the next row is NULL. The sqlmock suite below is still open and exercises a
realistic driver path.

Use `github.com/DATA-DOG/go-sqlmock`:

```go
db, mock, _ := sqlmock.New()
defer db.Close()

rows := sqlmock.NewRows([]string{"id", "name", "birthday", "active"}).
    AddRow(1, "Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), nil)
mock.ExpectQuery("SELECT …").WillReturnRows(rows)

var u User // Name types.NullString, Birthday types.NullDate, Active types.NullBool
_ = db.QueryRow("SELECT …").Scan(&u.ID, &u.Name, &u.Birthday, &u.Active)

// assert u.Active.Valid == false, u.Birthday.Valid == true, …
```

Cover: NULL in every column, valid values, and unsupported source type.

### 4.2. JSON round-trip through a realistic DTO

Add an integration-style test struct that exercises **every** nullable
type alongside regular `time.Time`, `uuid.UUID`, `decimal.Decimal`, and
verify:

```
struct → json.Marshal → json.Unmarshal → DeepEqual
```

Separately, test the case where some fields are `null` in JSON and are
omitted in the struct.

## 5. Fuzzing

Files `types/fuzz_parse_test.go`, `types/fuzz_timezone_test.go`:

```go
func FuzzParseDateFromString(f *testing.F) {
    f.Add("2023-05-20")
    f.Add("20.05.2023")
    f.Fuzz(func(t *testing.T, s string) {
        _, _ = ParseDateFromString(s)
    })
}
```

Goals: parsers never panic on arbitrary input; `ParseTimezoneExtended` does
not slice out of bounds on short strings.

## 6. Benchmarks

`types/bench_test.go`:
- `BenchmarkNullDate_MarshalJSON`,
- `BenchmarkNullDate_UnmarshalJSON`,
- `BenchmarkParseDateTimeFromString`,
- `BenchmarkAssembleNullDateTimeTZ`.

Purpose is regression tracking; no absolute threshold is set.

## 7. Infrastructure

### 7.1. Makefile / Taskfile

```
test:
	go test -race -count=1 ./...

cover:
	go test -race -coverprofile=cover.out ./...
	go tool cover -func=cover.out
	go tool cover -html=cover.out -o cover.html

lint:
	golangci-lint run
```

### 7.2. GitHub Actions (`.github/workflows/ci.yml`)

- matrix: Go 1.22.x, 1.23.x, 1.24.x, 1.26.x,
- `go test -race -coverprofile`,
- upload to Codecov,
- `golangci-lint` (gofmt, govet, staticcheck, errcheck, unused,
  goimports),
- fail PRs that drop coverage below the baseline.

### 7.3. Local pre-commit

`lefthook` or `pre-commit` running `go vet ./... && gofmt -l . && go test
./...`.

## 8. Phased rollout

| Phase | Scope                                                                 | Size |
| ----- | --------------------------------------------------------------------- | :--: |
| T1    | Table-driven unit tests for existing nullable types                   |  S   |
| T2    | DB-NULL edge cases in Scan (via `sqlmock`)                            |  M   |
| T3    | JSON round-trip for every type, including `NullUuid`, `NullDecimal`   |  S   |
| T4    | Tests for `emptiable`, `types_utils`, `string_able`                   |  S   |
| T5    | Tests for the new non-null Date/Time/DateTime (after `IMPROVE-PLAN` M1) |  M   |
| T6    | Tests for `LocalDateTime` / `NullLocalDateTime` (after `IMPROVE-PLAN` M2) |  M   |
| T7    | Parser fuzzing                                                        |  S   |
| T8    | Benchmarks with baseline numbers in the README                        |  S   |
| T9    | CI pipeline + Codecov / badges                                        |  S   |

Sizes: S ≈ up to one working day, M ≈ two–three days.