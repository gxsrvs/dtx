// Package types provides nullable wrappers and not-null helpers for values
// that must travel through both database/sql and JSON.
//
// Every wrapper implements sql.Scanner and driver.Valuer so the types can
// be used directly with sql.Row.Scan and db.Exec, and every wrapper
// implements json.Marshaler / json.Unmarshaler with a canonical shape: an
// invalid value marshals to null, a valid one to its documented string
// form (for example "2006-01-02" for Date/NullDate). Each type exposes an
// IsEmpty method (see [Emptiable]) and a ToString method (see
// [ToStringAble]) that returns "" when the value is NULL.
//
// The package is organised around pairs that share a single
// serializer/parser:
//
//   - Date / NullDate — calendar date ("2006-01-02" or "02.01.2006" on
//     input; always "2006-01-02" on output).
//   - LocalDateTime / NullLocalDateTime — wall-clock date-time without a
//     timezone ("2006-01-02T15:04:05[.fff…]").
//   - OffsetDateTime / NullOffsetDateTime — date-time with an offset
//     ("2006-01-02T15:04:05[.fff…]Z" or "±HH:MM").
//   - LocalTime / NullLocalTime — wall-clock time of day without a
//     timezone ("15:04:05[.fff…]").
//   - OffsetTime / NullOffsetTime — time of day with an offset
//     ("15:04:05[.fff…]Z" or "±HH:MM").
//
// Additional nullable wrappers cover primitive types: NullString,
// NullBool, NullInt16/32/64, NullFloat, NullDecimal and NullUuid.
//
// Nullable types accept the JSON token null on UnmarshalJSON and SQL NULL
// on Scan; their not-null counterparts reject both. Scan delegates to the
// corresponding sql.Null* wrapper so that typed nils from database
// drivers are honoured.
package types
