package repository

import (
	"aula4/internal/repository/storage"
	"errors"
)

type MockRepository struct {
	Products map[string]*storage.Product
}

func NewRepositoryProductsMock() MockRepository {
	return MockRepository{
		Products: make(map[string]*storage.Product),
	}
}

func (m *MockRepository) GetById(id string) (*storage.Product, error) {
	if product, exists := m.Products[id]; exists {
		return product, nil
	}
	return nil, errors.New("product not found")
}

func (m *MockRepository) GetAll() ([]*storage.Product, error) {
	var Products []*storage.Product
	for _, product := range m.Products {
		Products = append(Products, product)
	}

	if len(Products) == 0 {
		return nil, errors.New("no Products")
	}

	return Products, nil
}

func (m *MockRepository) Create(product storage.Product) (storage.Product, error) {
	m.Products[product.Id] = &product
	return product, nil
}

func (m *MockRepository) Update(product storage.Product) (storage.Product, error) {
	if _, exists := m.Products[product.Id]; exists {
		m.Products[product.Id] = &product
		return product, nil
	}

	return storage.Product{}, errors.New("product not found")
}

func (m *MockRepository) Patch(id string, updates map[string]interface{}) (*storage.Product, error) {
	if product, exists := m.Products[id]; exists {
		if name, ok := updates["name"].(string); ok {
			product.Name = name
		}
		if quantity, ok := updates["quantity"].(int); ok {
			product.Quantity = quantity
		}
		if codeValue, ok := updates["code_value"].(string); ok {
			product.Code_value = codeValue
		}
		if isPublished, ok := updates["is_published"].(bool); ok {
			product.Is_published = &isPublished
		}
		if expiration, ok := updates["expiration"].(string); ok {
			product.Expiration = expiration
		}
		if price, ok := updates["price"].(float64); ok {
			product.Price = price
		}

		return product, nil
	}

	return nil, errors.New("product not found")
}

func (m *MockRepository) Delete(id string) error {
	if _, exists := m.Products[id]; exists {
		delete(m.Products, id)
		return nil
	}
	return errors.New("product not found")
}
