// Package utils contains small JSON helpers used alongside the nullable
// and not-null types from package types.
//
// [LoadObjectFromJSON] / [LoadCollectionFromJSON] parse a JSON string
// into a strongly-typed value; the *FromJSONFile variants read the
// payload from disk first. [ToJSON] is the inverse — it serializes an
// arbitrary value and returns both the JSON string and any
// encoding/json error so callers can react to marshalling failures.
package utils
