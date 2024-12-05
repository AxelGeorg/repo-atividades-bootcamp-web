package model

import (
	"aula4/internal/repository"
	"aula4/internal/service/storage"
)

type ServiceProducts struct {
	Storage repository.RepositoryDB
}

func NewServiceProducts(storage repository.RepositoryDB) *ServiceProducts {
	return &ServiceProducts{
		Storage: storage,
	}
}

func (s *ServiceProducts) Create(product storage.Product) (storage.Product, error) {
	product, err := s.Storage.Create(product)
	if err != nil {
		return storage.Product{}, err
	}

	return product, nil
}

func (s *ServiceProducts) Update(product storage.Product) (storage.Product, error) {
	product, err := s.Storage.Update(product)
	if err != nil {
		return storage.Product{}, err
	}

	return product, nil
}
