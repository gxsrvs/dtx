package types

// nullJSON is the canonical JSON null token reused by every nullable
// type's MarshalJSON implementation, so that the binary representation
// stays identical across types.
var nullJSON = []byte("null")

// jsonTrue and jsonFalse are shared literal slices used by NullBool's
// MarshalJSON to skip allocation on the valid path. Treat them as
// immutable.
var (
	jsonTrue  = []byte("true")
	jsonFalse = []byte("false")
)
