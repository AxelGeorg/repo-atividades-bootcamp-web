package service

import (
	"aula4/internal/repository"
	"aula4/internal/repository/storage"
	"aula4/internal/validation"
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

func (s *ServiceProducts) GetTotalPrice(ids []string) (float64, []*storage.Product, error) {
	var (
		quantity int
		products []*storage.Product
		err      error
	)

	if len(ids) != 0 {
		quantity, products, err = GetProductQuantity(s, ids)
	} else {
		quantity, products, err = GetProductQuantityTotal(s)
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
	if err := validation.ValidateRequiredFields(product); err != nil {
		return storage.Product{}, err
	}

	if err := validation.CheckUniqueCodeValue(s.Repository, product.Code_value); err != nil {
		return storage.Product{}, err
	}

	if err := validation.ValidateDate(product.Expiration); err != nil {
		return storage.Product{}, err
	}

	product, err := s.Repository.Create(product)
	if err != nil {
		return storage.Product{}, err
	}

	return product, nil
}

func (s *ServiceProducts) Update(product storage.Product) (storage.Product, error) {
	if err := validation.ValidateRequiredFields(product); err != nil {
		return storage.Product{}, err
	}

	if err := validation.CheckUniqueCodeValue(s.Repository, product.Code_value); err != nil {
		return storage.Product{}, err
	}

	if err := validation.ValidateDate(product.Expiration); err != nil {
		return storage.Product{}, err
	}

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
