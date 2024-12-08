package service

import "aula4/internal/repository/storage"

type Service interface {
	GetAll() ([]*storage.Product, error)
	GetById(id string) (*storage.Product, error)
	Create(product storage.Product) (storage.Product, error)
	Update(product storage.Product) (storage.Product, error)
	Patch(id string, updates map[string]interface{}) (*storage.Product, error)
	Delete(id string) error

	SearchByPrice(price float64) ([]*storage.Product, error)
	GetProductQuantity(ids []string) (int, []*storage.Product, error)
	GetProductQuantityTotal() (int, []*storage.Product, error)
	GetTotalPrice(ids []string) (float64, []*storage.Product, error)
}
