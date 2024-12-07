package model

import (
	"aula4/internal/repository"
	"aula4/internal/service/storage"
	"errors"
	"log"
	"strings"
)

const (
	TaxLessThanTen         = 1.21
	TaxBetweenTenAndTwenty = 1.17
	TaxGreaterThanTwenty   = 1.15
)

const (
	countProductMin = 10
	countProductMax = 20
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

func (s *ServiceProducts) getProductQuantity(ids []string) (int, []*storage.Product, error) {
	mapProdQtd := make(map[string]int)

	var quantity int
	var products []*storage.Product
	for _, id := range ids {
		product, err := s.GetById(strings.TrimSpace(id))
		if err != nil {
			log.Printf("Error retrieving product with ID %s: %v", id, err)
			continue
		}

		mapProdQtd[id]++
		if mapProdQtd[id] > product.Quantity {
			return 0.0, nil, errors.New("not enough stock for product ID %s")
		}

		quantity++

		s.Update(*product)
		products = append(products, product)
	}

	if len(products) == 0 {
		return 0.0, nil, errors.New("no products available")
	}

	return quantity, products, nil
}

func (s *ServiceProducts) getProductQuantityTotal() (int, []*storage.Product, error) {
	products, err := s.GetAll()
	if err != nil {
		return 0.0, nil, err
	}

	var quantity int
	for _, product := range products {
		quantity += product.Quantity
	}

	return quantity, products, nil
}

func (s *ServiceProducts) TotalPrice(ids []string) (float64, []*storage.Product, error) {
	var (
		quantity int
		products []*storage.Product
		err      error
	)

	if len(ids) != 0 {
		quantity, products, err = s.getProductQuantity(ids)
	} else {
		quantity, products, err = s.getProductQuantityTotal()
	}

	if err != nil {
		return 0.0, nil, err
	}

	var totalPrice float64
	for _, product := range products {
		totalPrice += product.Price
	}

	var tax float64
	if quantity < countProductMin {
		tax = TaxLessThanTen
	} else if quantity >= countProductMin && quantity < countProductMax {
		tax = TaxBetweenTenAndTwenty
	} else {
		tax = TaxGreaterThanTwenty
	}

	totalPrice = totalPrice * tax
	return totalPrice, products, nil
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
