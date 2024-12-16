package store

type DB[T any] interface {
	Add(item T) (int, error)
	GetAll() ([]T, error)
	Change(id int, item T) error
	GetByID(id int) (T, error)
	GetByField(field string, value any) ([]T, error)
	Remove(id int) error
}
