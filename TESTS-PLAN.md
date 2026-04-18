# `gxsrvs/dtx` — Testing Plan

The goal is to bring library test coverage up to public-release quality
(target ≥ 85% line coverage; JSON-null / SQL-NULL / invalid-input branches
are mandatory) and to formalise the coverage gate in CI.

Tests are plain `testing` + table-driven, without third-party frameworks.
Integration scenarios that require `database/sql` semantics use
`github.com/DATA-DOG/go-sqlmock`.

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

### 2.1. `types/null_string.go` — `null_string_test.go`
Already has baseline tests. Extend with:
- UnmarshalJSON for escaped characters (`"a\"b"`, `"\u0026"`),
- Scan from `[]byte`, from `int` (error), from `time.Time` (error),
- round-trip for non-ASCII strings,
- `IsEmpty` reached through the `Emptiable` interface.

### 2.2. `types/null_bool.go` — `null_bool_test.go` *(new)*
- `NewNullBool(true)`, `NewNullBoolEmpty()`, `NullBoolFromString`,
- Scan from `bool`, `int64` (1/0), `[]byte("t")`, `nil`,
- UnmarshalJSON: `true`, `false`, `null`, `"true"` (stringified),
- MarshalJSON: invalid → `null`.

### 2.3. `types/null_int16.go`, `null_int32.go`, `null_int64.go`
Partial coverage exists. Extend with:
- boundary values (`math.MaxInt16`, `math.MinInt16`, overflow on
  `FromString`),
- Scan from `int64`, `int32`, `string("42")`, `nil`, and an unsupported
  type,
- UnmarshalJSON for numeric, stringified (`"42"`), `null`, floating
  (`42.0`).

### 2.4. `types/null_float.go` — `null_float_test.go` *(new)*
- Boundary values (`+Inf`, `-Inf`, `NaN`: decide in the API whether they
  are supported, then test).
- Verify that `ToString()` returns `""` once the inconsistency documented
  in `IMPROVE-PLAN.md` (M4) is fixed.

### 2.5. `types/null_decimal.go` — `null_decimal_test.go`
Extend with:
- JSON round-trip on fractional values (0.1 + 0.2),
- Scan from `string`, `[]byte`, `nil`; from `float64` (error),
- `MulNullDecimals` — valid/invalid boundary,
- `MarshalJSON` emits a bare number (no quotes) — check JS/TS
  compatibility.

### 2.6. `types/null_uuid.go`
Extend with:
- Scan from a UUID string, from `[]byte`, from `nil`,
- `MarshalJSON` → `"xxxxxxxx-xxxx-…"`,
- `UnmarshalJSON` for canonical, braced, upper-case, `null`.

### 2.7. `types/null_date.go`, `null_iso_date.go`
Partial coverage exists. Extend with:
- `ParseDateFromString` for both formats (`2006-01-02`, `02.01.2006`),
- invalid inputs: `2023-13-01`, `31.02.2023`, `0000-00-00`, `2023/01/02`,
- Scan from `time.Time`, `string`, `nil`,
- UnmarshalJSON for `"2023-05-20"`, `"20.05.2023"`, `null`,
- `DateToString` across time zones.

### 2.8. `types/null_datetime.go`, `datetime.go`
- parse both separators (`T` and space),
- fractional seconds (`2023-05-20T15:30:45.123`),
- with timezone (`2023-05-20T15:30:45+03:00`,
  `2023-05-20 15:30:45+0300`),
- Scan from `time.Time`, `string`, `[]byte`, `nil`,
- for the non-null `DateTime`: cover empty `Scan` branches, add `Value()`
  (see `IMPROVE-PLAN.md` M1) and test it.

### 2.9. `types/null_time.go`, `time_only.go`
- formats `HH:MM`, `HH:MM:SS`, `HH:MM:SS.ms`,
- with timezone `15:30:00+0300`, `15:30:00+03:00`,
- direct unit tests for `isDigitString`, `parseTimezone`,
  `parseTimezoneWithColon`,
- Scan / Value for `TimeOnly`,
- clarify the semantic gap between `NullTime` and the new
  `NullLocalDateTime` (see `IMPROVE-PLAN.md` M2).

### 2.10. `types/emptiable.go`, `types_utils.go`, `string_able.go`
- `IsEmpty` for every library type (both value and pointer),
- `ToString` for every type plus the default branch,
- `MaxDateTime`, `MinDateTime`,
- `AssembleDateTime` with `nil` location and with a passed-in one,
- `AssembleDateTimeTZ` — correct zone and invalid zone string,
- `AssembleNullDateTimeTZ` — all four valid/invalid combinations for date
  and time.

### 2.11. `utils/json.go`
Baseline exists. Extend with:
- `LoadObjectFromJson` with a generic type that contains
  `NullString`/`NullDate` — confirm round-trip preservation,
- `LoadCollectionFromJsonFile` — empty file, non-array input (error),
- `ToJson` — once the API returns an error (see `IMPROVE-PLAN.md` M8),
  verify the returned error; until then, assert `""` on failure.

### 2.12. `dto/data.go`
Baseline exists. Extend with:
- JSON round-trip for `StdDataPackage[T]` using library types
  (`StdDataPackage[UserDTO]`, where `UserDTO.Birthday types.NullDate`).

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

In `dto/`, add a struct that exercises **every** nullable type alongside
regular `time.Time`, `uuid.UUID`, `decimal.Decimal`, and verify:

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