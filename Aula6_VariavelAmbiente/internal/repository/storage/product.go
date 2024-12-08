package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
)

const (
	localFileJson = "../../docs/db/json/products.json"
)

type Product struct {
	Id           string
	Name         string
	Quantity     int
	Code_value   string
	Is_published bool
	Expiration   string
	Price        float64
}

type StorageProducts struct {
	mu sync.Mutex
}

func NewStorageProducts() StorageProducts {
	return StorageProducts{}
}

func (s *StorageProducts) ReadAllProductsToFile() ([]*Product, error) {
	var productList []*Product

	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(localFileJson)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(localFileJson)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			initialData := []Product{}
			writer := json.NewEncoder(file)
			if err := writer.Encode(initialData); err != nil {
				return nil, err
			}

			return productList, nil
		}
		return nil, err
	}
	defer file.Close()

	reader := json.NewDecoder(file)
	err = reader.Decode(&productList)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	return productList, nil
}

func (s *StorageProducts) ReadProductById(id string) (*Product, error) {
	products, err := s.ReadAllProductsToFile()
	if err != nil {
		return nil, err
	}
	for _, product := range products {
		if product.Id == id {
			return product, nil
		}
	}
	return nil, nil
}

func (s *StorageProducts) SaveProduct(product *Product) error {
	products, err := s.ReadAllProductsToFile()
	if err != nil {
		return err
	}

	for _, p := range products {
		if p.Id == product.Id {
			return errors.New("product already exists")
		}
	}

	products = append(products, product)
	return s.WriteProductsToFile(products)
}

func (s *StorageProducts) UpdateProduct(updatedProduct *Product) error {
	products, err := s.ReadAllProductsToFile()
	if err != nil {
		return err
	}

	for i, product := range products {
		if product.Id == updatedProduct.Id {
			products[i] = updatedProduct
			return s.WriteProductsToFile(products)
		}
	}
	return errors.New("product not found")
}

func (s *StorageProducts) DeleteProduct(id string) error {
	products, err := s.ReadAllProductsToFile()
	if err != nil {
		return err
	}

	for i, product := range products {
		if product.Id == id {
			products = append(products[:i], products[i+1:]...)
			return s.WriteProductsToFile(products)
		}
	}

	return errors.New("product not found")
}

func (s *StorageProducts) WriteProductsToFile(productList []*Product) error {
	file, err := os.OpenFile(localFileJson, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := json.NewEncoder(file)
	return writer.Encode(productList)
}
