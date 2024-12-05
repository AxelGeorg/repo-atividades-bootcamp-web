package repository

import (
	internal "aula4/internal/storage"

	"github.com/google/uuid"
)

type RepositoryDB struct {
	DB map[string]*internal.Product
}

func (r *RepositoryDB) Create(product internal.Product) (internal.Product, error) {
	id := uuid.New()
	product.Id = id.String()

	r.DB[id.String()] = &product
	return product, nil
}

func NewMeliDB() RepositoryDB {
	return RepositoryDB{
		DB: make(map[string]*internal.Product),
	}
}
