package repository

import "aula4/internal/repository/storage"

type Repository interface {
	GetById(id string) (*storage.Product, error)
	GetAll() ([]*storage.Product, error)
	Create(product storage.Product) (storage.Product, error)
	Update(product storage.Product) (storage.Product, error)
	Patch(id string, updates map[string]interface{}) (*storage.Product, error)
	Delete(id string) error
}
