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

func TestLoadObjectFromJSON(t *testing.T) {
	jsonStr := `{"id": 1, "name": "Test"}`
	obj, err := LoadObjectFromJSON[testObject](jsonStr)
	if err != nil {
		t.Fatalf("LoadObjectFromJSON failed: %v", err)
	}

	expected := &testObject{ID: 1, Name: "Test"}
	if !reflect.DeepEqual(obj, expected) {
		t.Errorf("got %+v, want %+v", obj, expected)
	}

	// Test error case
	_, err = LoadObjectFromJSON[testObject](`invalid json`)
	if err == nil {
		t.Error("expected error for invalid json, got nil")
	}
}

func TestLoadCollectionFromJSON(t *testing.T) {
	jsonStr := `[{"id": 1, "name": "Test1"}, {"id": 2, "name": "Test2"}]`
	coll, err := LoadCollectionFromJSON[testObject](jsonStr)
	if err != nil {
		t.Fatalf("LoadCollectionFromJSON failed: %v", err)
	}

	expected := []testObject{{ID: 1, Name: "Test1"}, {ID: 2, Name: "Test2"}}
	if !reflect.DeepEqual(coll, expected) {
		t.Errorf("got %+v, want %+v", coll, expected)
	}

	// Test error case
	_, err = LoadCollectionFromJSON[testObject](`invalid json`)
	if err == nil {
		t.Error("expected error for invalid json, got nil")
	}
}

func TestLoadObjectFromJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	fileName := filepath.Join(tmpDir, "test_obj.json")
	jsonStr := `{"id": 10, "name": "FileTest"}`

	if err := os.WriteFile(fileName, []byte(jsonStr), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	obj, err := LoadObjectFromJSONFile[testObject](fileName)
	if err != nil {
		t.Fatalf("LoadObjectFromJSONFile failed: %v", err)
	}

	expected := &testObject{ID: 10, Name: "FileTest"}
	if !reflect.DeepEqual(obj, expected) {
		t.Errorf("got %+v, want %+v", obj, expected)
	}

	// Test file not found
	_, err = LoadObjectFromJSONFile[testObject]("non_existent.json")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}

	// Test invalid json in file
	invalidFileName := filepath.Join(tmpDir, "invalid.json")
	_ = os.WriteFile(invalidFileName, []byte(`{invalid}`), 0644)
	_, err = LoadObjectFromJSONFile[testObject](invalidFileName)
	if err == nil {
		t.Error("expected error for invalid json in file, got nil")
	}
}

func TestLoadCollectionFromJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	fileName := filepath.Join(tmpDir, "test_coll.json")
	jsonStr := `[{"id": 1, "name": "A"}, {"id": 2, "name": "B"}]`

	if err := os.WriteFile(fileName, []byte(jsonStr), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	coll, err := LoadCollectionFromJSONFile[testObject](fileName)
	if err != nil {
		t.Fatalf("LoadCollectionFromJSONFile failed: %v", err)
	}

	expected := []testObject{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}
	if !reflect.DeepEqual(coll, expected) {
		t.Errorf("got %+v, want %+v", coll, expected)
	}

	// Test file not found
	_, err = LoadCollectionFromJSONFile[testObject]("non_existent.json")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestToJSON(t *testing.T) {
	obj := testObject{ID: 5, Name: "JSON"}
	expected := `{"id":5,"name":"JSON"}`
	result, err := ToJSON(obj)
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}
	if result != expected {
		t.Errorf("got %s, want %s", result, expected)
	}

	// Test map
	m := map[string]string{"key": "value"}
	expectedMap := `{"key":"value"}`
	resultMap, err := ToJSON(m)
	if err != nil {
		t.Fatalf("ToJSON(map) failed: %v", err)
	}
	if resultMap != expectedMap {
		t.Errorf("got %s, want %s", resultMap, expectedMap)
	}

	// Channels cannot be marshaled to JSON — expect an error.
	ch := make(chan int)
	result, err = ToJSON(ch)
	if err == nil {
		t.Error("expected error for unmarshalable value, got nil")
	}
	if result != "" {
		t.Errorf("expected empty string on error, got %s", result)
	}
}
