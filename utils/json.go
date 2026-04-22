package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadObjectFromJSON unmarshals a JSON object into a fresh T and returns
// a pointer to it. Any json.Unmarshal error is returned as-is.
func LoadObjectFromJSON[T any](jsonString string) (*T, error) {
	var item T
	if err := json.Unmarshal([]byte(jsonString), &item); err != nil {
		return nil, err
	}
	return &item, nil
}

// LoadCollectionFromJSON unmarshals a JSON array into a []T. Any
// json.Unmarshal error is returned as-is.
func LoadCollectionFromJSON[T any](jsonString string) ([]T, error) {
	var result []T
	if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// LoadObjectFromJSONFile reads fileName from disk and unmarshals the
// contents into a fresh T. Read and unmarshal errors are wrapped with
// the file name for easier diagnosis.
func LoadObjectFromJSONFile[T any](fileName string) (*T, error) {
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

// LoadCollectionFromJSONFile reads fileName from disk and unmarshals the
// contents into a []T. Read and unmarshal errors are wrapped with the
// file name for easier diagnosis.
func LoadCollectionFromJSONFile[T any](fileName string) ([]T, error) {
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

// ToJSON marshals entity to a JSON string. Marshal errors are wrapped
// so callers can distinguish them from downstream failures.
func ToJSON(entity interface{}) (string, error) {
	b, err := json.Marshal(entity)
	if err != nil {
		return "", fmt.Errorf("failed to marshal interface to json: %w", err)
	}
	return string(b), nil
}
