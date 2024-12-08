package repository

import (
	"aula4/internal/repository/storage"
	"errors"

	"github.com/google/uuid"
)

type RepositoryProducts struct {
	Storage storage.Storage
}

func NewRepositoryProducts(storage storage.Storage) RepositoryProducts {
	return RepositoryProducts{
		Storage: storage,
	}
}

func (r *RepositoryProducts) GetById(id string) (*storage.Product, error) {
	product, err := r.Storage.ReadProductById(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (r *RepositoryProducts) GetAll() ([]*storage.Product, error) {
	products, err := r.Storage.ReadAllProductsToFile()
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, errors.New("no products")
	}

	return products, nil
}

func (r *RepositoryProducts) Create(product storage.Product) (storage.Product, error) {
	id := uuid.New()
	product.Id = id.String()

	if err := r.Storage.SaveProduct(&product); err != nil {
		return storage.Product{}, err
	}

	return product, nil
}

func (r *RepositoryProducts) Update(product storage.Product) (storage.Product, error) {
	if err := r.Storage.UpdateProduct(&product); err != nil {
		return storage.Product{}, err
	}

	return product, nil
}

func (r *RepositoryProducts) Patch(id string, updates map[string]interface{}) (*storage.Product, error) {
	product, err := r.GetById(id)
	if err != nil {
		return nil, err
	}

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
		product.Is_published = isPublished
	}
	if expiration, ok := updates["expiration"].(string); ok {
		product.Expiration = expiration
	}
	if price, ok := updates["price"].(float64); ok {
		product.Price = price
	}

	if err := r.Storage.UpdateProduct(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (r *RepositoryProducts) Delete(id string) error {
	return r.Storage.DeleteProduct(id)
}
