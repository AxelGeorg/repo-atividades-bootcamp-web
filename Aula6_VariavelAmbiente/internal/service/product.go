package service

import (
	"aula4/internal/repository"
	"aula4/internal/repository/storage"
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
	Repository repository.Repository
}

func NewServiceProducts(repository repository.Repository) ServiceProducts {
	return ServiceProducts{
		Repository: repository,
	}
}

func (s *ServiceProducts) SearchByPrice(price float64) ([]*storage.Product, error) {
	products, err := s.Repository.GetAll()
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

func (s *ServiceProducts) GetProductQuantity(ids []string) (int, []*storage.Product, error) {
	mapProdQtd := make(map[string]int)

	var quantity int
	var products []*storage.Product
	for _, id := range ids {
		idStr := strings.TrimSpace(id)
		product, err := s.GetById(idStr)
		if err != nil {
			log.Printf("Error retrieving product with ID:" + idStr + " - Error: " + err.Error())
			continue
		}

		mapProdQtd[idStr]++
		if mapProdQtd[idStr] > product.Quantity {
			return 0.0, nil, errors.New("not enough stock for product ID:" + idStr)
		}

		product.Is_published = true
		_, err = s.Update(*product)
		if err != nil {
			return 0.0, nil, err
		}

		quantity++
		products = append(products, product)
	}

	if len(products) == 0 {
		return 0.0, nil, errors.New("no products available")
	}

	return quantity, products, nil
}

func (s *ServiceProducts) GetProductQuantityTotal() (int, []*storage.Product, error) {
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

func (s *ServiceProducts) GetTotalPrice(ids []string) (float64, []*storage.Product, error) {
	var (
		quantity int
		products []*storage.Product
		err      error
	)

	if len(ids) != 0 {
		quantity, products, err = s.GetProductQuantity(ids)
	} else {
		quantity, products, err = s.GetProductQuantityTotal()
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
	products, err := s.Repository.GetAll()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ServiceProducts) GetById(id string) (*storage.Product, error) {
	product, err := s.Repository.GetById(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ServiceProducts) Create(product storage.Product) (storage.Product, error) {
	product, err := s.Repository.Create(product)
	if err != nil {
		return storage.Product{}, err
	}

	return product, nil
}

func (s *ServiceProducts) Update(product storage.Product) (storage.Product, error) {
	product, err := s.Repository.Update(product)
	if err != nil {
		return storage.Product{}, err
	}

	return product, nil
}

func (s *ServiceProducts) Patch(id string, updates map[string]interface{}) (*storage.Product, error) {
	product, err := s.Repository.Patch(id, updates)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ServiceProducts) Delete(id string) error {
	err := s.Repository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
