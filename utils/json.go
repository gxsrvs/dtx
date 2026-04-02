package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func LoadObjectFromJson[T any](jsonString string) (*T, error) {
	var item T
	if err := json.Unmarshal([]byte(jsonString), &item); err != nil {
		return nil, err
	}
	return &item, nil
}

//goland:noinspection GoUnusedExportedFunction
func LoadCollectionFromJson[T any](jsonString string) ([]T, error) {
	var result []T
	if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
		return nil, err
	}
	return result, nil
}

//goland:noinspection GoUnusedExportedFunction
func LoadObjectFromJsonFile[T any](fileName string) (*T, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fileName, err)
	}

	var item T
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data from %s: %w", fileName, err)
	}
	return &item, nil
}

func LoadCollectionFromJsonFile[T any](fileName string) ([]T, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fileName, err)
	}

	var result []T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data from %s: %w", fileName, err)
	}
	return result, nil
}

func ToJson(entity interface{}) string {
	b, err := json.Marshal(entity)
	if err != nil {
		log.Println(fmt.Errorf("failed to marshal interface to json: %w", err))
		return ""
	}
	return string(b)
}
