package validation

import (
	"aula4/internal/repository"
	"aula4/internal/repository/storage"
	"errors"
	"time"
)

func CheckUniqueCodeValue(repo repository.Repository, codeValue string) error {
	products, err := repo.GetAll()
	if err != nil {
		if err.Error() != "no Products" {
			return err
		}

		return nil
	}

	for _, product := range products {
		if product.Code_value == codeValue {
			return errors.New("the code_value must be unique")
		}
	}

	return nil
}

func ValidateDate(dateStr string) error {
	layout := "02/02/2006"
	_, err := time.Parse(layout, dateStr)
	if err != nil {
		return errors.New("data inválida. O formato deve ser DD/MM/YYYY")
	}
	return nil
}

func ValidateRequiredFields(product storage.Product) error {
	if product.Code_value == "" || product.Name == "" || product.Quantity <= 0 ||
		product.Expiration == "" || product.Price <= 0 {
		return errors.New("all fields must be filled, except for is_published")
	}
	return nil
}
