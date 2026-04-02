package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type testObject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestLoadObjectFromJson(t *testing.T) {
	jsonStr := `{"id": 1, "name": "Test"}`
	obj, err := LoadObjectFromJson[testObject](jsonStr)
	if err != nil {
		t.Fatalf("LoadObjectFromJson failed: %v", err)
	}

	expected := &testObject{ID: 1, Name: "Test"}
	if !reflect.DeepEqual(obj, expected) {
		t.Errorf("got %+v, want %+v", obj, expected)
	}

	// Test error case
	_, err = LoadObjectFromJson[testObject](`invalid json`)
	if err == nil {
		t.Error("expected error for invalid json, got nil")
	}
}

func TestLoadCollectionFromJson(t *testing.T) {
	jsonStr := `[{"id": 1, "name": "Test1"}, {"id": 2, "name": "Test2"}]`
	coll, err := LoadCollectionFromJson[testObject](jsonStr)
	if err != nil {
		t.Fatalf("LoadCollectionFromJson failed: %v", err)
	}

	expected := []testObject{{ID: 1, Name: "Test1"}, {ID: 2, Name: "Test2"}}
	if !reflect.DeepEqual(coll, expected) {
		t.Errorf("got %+v, want %+v", coll, expected)
	}

	// Test error case
	_, err = LoadCollectionFromJson[testObject](`invalid json`)
	if err == nil {
		t.Error("expected error for invalid json, got nil")
	}
}

func TestLoadObjectFromJsonFile(t *testing.T) {
	tmpDir := t.TempDir()
	fileName := filepath.Join(tmpDir, "test_obj.json")
	jsonStr := `{"id": 10, "name": "FileTest"}`

	if err := os.WriteFile(fileName, []byte(jsonStr), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	obj, err := LoadObjectFromJsonFile[testObject](fileName)
	if err != nil {
		t.Fatalf("LoadObjectFromJsonFile failed: %v", err)
	}

	expected := &testObject{ID: 10, Name: "FileTest"}
	if !reflect.DeepEqual(obj, expected) {
		t.Errorf("got %+v, want %+v", obj, expected)
	}

	// Test file not found
	_, err = LoadObjectFromJsonFile[testObject]("non_existent.json")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}

	// Test invalid json in file
	invalidFileName := filepath.Join(tmpDir, "invalid.json")
	_ = os.WriteFile(invalidFileName, []byte(`{invalid}`), 0644)
	_, err = LoadObjectFromJsonFile[testObject](invalidFileName)
	if err == nil {
		t.Error("expected error for invalid json in file, got nil")
	}
}

func TestLoadCollectionFromJsonFile(t *testing.T) {
	tmpDir := t.TempDir()
	fileName := filepath.Join(tmpDir, "test_coll.json")
	jsonStr := `[{"id": 1, "name": "A"}, {"id": 2, "name": "B"}]`

	if err := os.WriteFile(fileName, []byte(jsonStr), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	coll, err := LoadCollectionFromJsonFile[testObject](fileName)
	if err != nil {
		t.Fatalf("LoadCollectionFromJsonFile failed: %v", err)
	}

	expected := []testObject{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}
	if !reflect.DeepEqual(coll, expected) {
		t.Errorf("got %+v, want %+v", coll, expected)
	}

	// Test file not found
	_, err = LoadCollectionFromJsonFile[testObject]("non_existent.json")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestToJson(t *testing.T) {
	obj := testObject{ID: 5, Name: "JSON"}
	expected := `{"id":5,"name":"JSON"}`
	result := ToJson(obj)

	if result != expected {
		t.Errorf("got %s, want %s", result, expected)
	}

	// Test map
	m := map[string]string{"key": "value"}
	expectedMap := `{"key":"value"}`
	resultMap := ToJson(m)
	if resultMap != expectedMap {
		t.Errorf("got %s, want %s", resultMap, expectedMap)
	}

	// Test marshaling error (channels cannot be marshaled to JSON)
	ch := make(chan int)
	resultErr := ToJson(ch)
	if resultErr != "" {
		t.Errorf("expected empty string for marshal error, got %s", resultErr)
	}
}
