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

func (s *ServiceProducts) SearchByPrice(price float64) ([]*storage.Product, error) {
	products, err := s.Storage.GetAll()
	if err != nil {
		return nil, err
	}

	var filteredProducts []*storage.Product

	for _, product := range products {
		if product.Price > price {
			filteredProducts = append(filteredProducts, product)
		}
	}

	return filteredProducts, nil
}

func (s *ServiceProducts) GetAll() ([]*storage.Product, error) {
	products, err := s.Storage.GetAll()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ServiceProducts) GetById(id string) (*storage.Product, error) {
	product, err := s.Storage.GetById(id)
	if err != nil {
		return nil, err
	}

	return product, nil
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

func (s *ServiceProducts) Patch(id string, updates map[string]interface{}) (*storage.Product, error) {
	product, err := s.Storage.Patch(id, updates)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ServiceProducts) Delete(id string) error {
	err := s.Storage.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
