package types

// nullJson is the canonical JSON null token reused by every nullable
// type's MarshalJSON implementation, so that the binary representation
// stays identical across types.
var nullJson = []byte("null")
