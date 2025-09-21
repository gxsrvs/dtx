package dto

// DataObject - интерфейс объектов данных
type DataObject interface{}

// DataPackage - интерфейс для пакетов данных, содержащих множество DataObject
type DataPackage[T DataObject] interface {
	GetData() []T
}

type StdDataPackage[T DataObject] struct {
	Objects []T `json:"objects" validate:"required"`
}

func (d *StdDataPackage[T]) GetData() []T {
	return d.Objects
}
