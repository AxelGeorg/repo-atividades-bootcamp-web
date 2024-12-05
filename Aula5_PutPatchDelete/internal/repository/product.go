package repository

import (
	"aula4/internal/service/storage"

	"github.com/google/uuid"
)

type RepositoryDB struct {
	DB map[string]*storage.Product
}

func (r *RepositoryDB) Create(product storage.Product) (storage.Product, error) {
	id := uuid.New()
	product.Id = id.String()

	r.DB[id.String()] = &product
	return product, nil
}

func (r *RepositoryDB) Update(product storage.Product) (storage.Product, error) {
	r.DB[product.Id] = &product
	return product, nil
}

func (r *RepositoryDB) Delete(id string) error {
	delete(r.DB, id)
	return nil
}

func NewMeliDB() RepositoryDB {
	return RepositoryDB{
		DB: make(map[string]*storage.Product),
	}
}
