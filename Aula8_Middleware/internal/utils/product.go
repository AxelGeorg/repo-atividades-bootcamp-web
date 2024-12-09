package utils

import (
	"aula4/internal/repository/storage"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	MessageProductCreated = "Product created"
	MessageProductUpdated = "Product updated"
	MessageProductDeleted = "Product deleted"
)

func CheckUniqueCodeValue(products []*storage.Product, prod storage.Product) error {
	for _, product := range products {
		if product.Id == prod.Id {
			continue
		}

		if product.Code_value == prod.Code_value {
			return errors.New("the code_value must be unique")
		}
	}

	return nil
}

func ValidateUUID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("invalid UUID format: %s", id)
	}
	return nil
}

func ValidateDate(dateStr string) error {
	layout := "02/01/2006"
	_, err := time.Parse(layout, dateStr)
	if err != nil {
		return errors.New("invalid date. The format must be DD/MM/YYYY")
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
