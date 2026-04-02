package types

import (
	"strings"

	"github.com/google/uuid"
)

type NullUuid uuid.NullUUID

func NewNullUuidEmpty() NullUuid {
	return (NullUuid)(uuid.NullUUID{
		UUID:  uuid.UUID{},
		Valid: false,
	})
}

func NewNullUuid(val uuid.UUID) NullUuid {
	return (NullUuid)(uuid.NullUUID{
		UUID:  val,
		Valid: true,
	})
}

//goland:noinspection GoUnusedExportedFunction
func NullUuidToString(val uuid.NullUUID) string {
	if !val.Valid {
		return "null"
	}
	return val.UUID.String()
}

func NullUuidFromString(strValue *string) NullUuid {
	if strValue == nil || *strValue == "" ||
		strings.ToLower(*strValue) == "null" ||
		strings.ToLower(*strValue) == "nil" {
		return NewNullUuidEmpty()
	}
	result, err := uuid.Parse(*strValue)
	if err != nil {
		return NewNullUuidEmpty()
	}
	return NewNullUuid(result)
}
