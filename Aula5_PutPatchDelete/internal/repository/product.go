package repository

import (
	"aula4/internal/service/storage"
	"errors"

	"github.com/google/uuid"
)

type RepositoryDB struct {
	DB map[string]*storage.Product
}

func NewMeliDB() RepositoryDB {
	return RepositoryDB{
		DB: make(map[string]*storage.Product),
	}
}

func (r *RepositoryDB) GetById(id string) (*storage.Product, error) {
	product, exists := r.DB[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (r *RepositoryDB) GetAll() ([]*storage.Product, error) {
	products := make([]*storage.Product, 0, len(r.DB))

	if len(r.DB) == 0 {
		return nil, errors.New("no products available")
	}

	for _, product := range r.DB {
		products = append(products, product)
	}
	return products, nil
}

func (r *RepositoryDB) Create(product storage.Product) (storage.Product, error) {
	id := uuid.New()
	product.Id = id.String()

	if _, exists := r.DB[product.Id]; exists {
		return storage.Product{}, errors.New("product already exists")
	}

	r.DB[product.Id] = &product
	return product, nil
}

func (r *RepositoryDB) Update(product storage.Product) (storage.Product, error) {
	if _, exists := r.DB[product.Id]; !exists {
		return storage.Product{}, errors.New("product not found")
	}

	r.DB[product.Id] = &product
	return product, nil
}

func (r *RepositoryDB) Delete(id string) error {
	if _, exists := r.DB[id]; !exists {
		return errors.New("product not found")
	}

	delete(r.DB, id)
	return nil
}
