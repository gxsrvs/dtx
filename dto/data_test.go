package dto

import (
	"testing"
)

// TestUser представляет тестовый объект данных
type TestUser struct {
	Id   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

// TestProduct представляет другой тестовый объект данных
type TestProduct struct {
	Id    int     `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

// TestStdDataPackage_BasicOperations тестирует основные операции пакета данных
func TestStdDataPackage_BasicOperations(t *testing.T) {
	// Тест пустого пакета
	emptyPackage := &StdDataPackage[TestUser]{
		Objects: []TestUser{},
	}

	data := emptyPackage.GetData()
	if len(data) != 0 {
		t.Errorf("Expected empty data package, got %d items", len(data))
	}

	// Тест пакета с данными
	users := []TestUser{
		{Id: 1, Name: "Alice"},
		{Id: 2, Name: "Bob"},
		{Id: 3, Name: "Charlie"},
	}

	userPackage := &StdDataPackage[TestUser]{
		Objects: users,
	}

	retrievedData := userPackage.GetData()
	if len(retrievedData) != 3 {
		t.Errorf("Expected 3 users, got %d", len(retrievedData))
	}

	// Проверяем данные
	if retrievedData[0].Id != 1 || retrievedData[0].Name != "Alice" {
		t.Errorf("Expected first user to be {1, Alice}, got {%d, %s}",
			retrievedData[0].Id, retrievedData[0].Name)
	}

	if retrievedData[1].Id != 2 || retrievedData[1].Name != "Bob" {
		t.Errorf("Expected second user to be {2, Bob}, got {%d, %s}",
			retrievedData[1].Id, retrievedData[1].Name)
	}

	if retrievedData[2].Id != 3 || retrievedData[2].Name != "Charlie" {
		t.Errorf("Expected third user to be {3, Charlie}, got {%d, %s}",
			retrievedData[2].Id, retrievedData[2].Name)
	}
}

// TestStdDataPackage_DifferentTypes тестирует работу с разными типами данных
func TestStdDataPackage_DifferentTypes(t *testing.T) {
	// Тест с продуктами
	products := []TestProduct{
		{Id: 101, Title: "Laptop", Price: 999.99},
		{Id: 102, Title: "Mouse", Price: 29.99},
	}

	productPackage := &StdDataPackage[TestProduct]{
		Objects: products,
	}

	retrievedProducts := productPackage.GetData()
	if len(retrievedProducts) != 2 {
		t.Errorf("Expected 2 products, got %d", len(retrievedProducts))
	}

	if retrievedProducts[0].Title != "Laptop" || retrievedProducts[0].Price != 999.99 {
		t.Errorf("Expected first product to be {Laptop, 999.99}, got {%s, %.2f}",
			retrievedProducts[0].Title, retrievedProducts[0].Price)
	}

	// Тест с простыми типами (string)
	stringPackage := &StdDataPackage[string]{
		Objects: []string{"hello", "world", "test"},
	}

	retrievedStrings := stringPackage.GetData()
	if len(retrievedStrings) != 3 {
		t.Errorf("Expected 3 strings, got %d", len(retrievedStrings))
	}

	expectedStrings := []string{"hello", "world", "test"}
	for i, str := range retrievedStrings {
		if str != expectedStrings[i] {
			t.Errorf("Expected string at index %d to be '%s', got '%s'", i, expectedStrings[i], str)
		}
	}
}

// TestStdDataPackage_NilObjects тестирует обработку nil объектов
func TestStdDataPackage_NilObjects(t *testing.T) {
	// Тест с nil slice
	nilPackage := &StdDataPackage[TestUser]{
		Objects: nil,
	}

	data := nilPackage.GetData()
	if data != nil {
		t.Errorf("Expected nil data for nil Objects, got %v", data)
	}
}

// TestStdDataPackage_Modification тестирует изменение возвращаемых данных
func TestStdDataPackage_Modification(t *testing.T) {
	originalUsers := []TestUser{
		{Id: 1, Name: "Alice"},
		{Id: 2, Name: "Bob"},
	}

	userPackage := &StdDataPackage[TestUser]{
		Objects: originalUsers,
	}

	// Получаем данные и модифицируем их
	retrievedData := userPackage.GetData()
	if len(retrievedData) != 2 {
		t.Errorf("Expected 2 users, got %d", len(retrievedData))
	}

	// Модифицируем полученные данные
	retrievedData[0].Name = "Modified Alice"

	// Проверяем, что оригинальные данные тоже изменились (так как возвращается slice, а не копия)
	originalData := userPackage.GetData()
	if originalData[0].Name != "Modified Alice" {
		t.Error("Expected modification to affect original data (slice is not copied)")
	}
}

// TestStdDataPackage_Interface тестирует соответствие интерфейсу
func TestStdDataPackage_Interface(t *testing.T) {
	// Проверяем, что StdDataPackage реализует интерфейс DataPackage
	var _ DataPackage[TestUser] = &StdDataPackage[TestUser]{}
	var _ DataPackage[string] = &StdDataPackage[string]{}
	var _ DataPackage[int] = &StdDataPackage[int]{}

	// Тест через интерфейс
	var userPackageInterface DataPackage[TestUser] = &StdDataPackage[TestUser]{
		Objects: []TestUser{{Id: 1, Name: "Interface Test"}},
	}

	data := userPackageInterface.GetData()
	if len(data) != 1 {
		t.Errorf("Expected 1 user through interface, got %d", len(data))
	}

	if data[0].Name != "Interface Test" {
		t.Errorf("Expected user name 'Interface Test', got '%s'", data[0].Name)
	}
}
