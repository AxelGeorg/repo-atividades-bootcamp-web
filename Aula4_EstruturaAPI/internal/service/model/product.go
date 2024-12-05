package model

import (
	"aula4/internal"
	"aula4/internal/repository"
)

type ServiceProducts struct {
	Storage repository.RepositoryDB
}

func NewServiceProducts(storage repository.RepositoryDB) *ServiceProducts {
	return &ServiceProducts{
		Storage: storage,
	}
}

func (s *ServiceProducts) Create(product internal.Product) (internal.Product, error) {
	product, err := s.Storage.Create(product)
	if err != nil {
		return internal.Product{}, err
	}

	return product, nil
}
