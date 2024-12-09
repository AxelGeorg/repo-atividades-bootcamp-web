package utils

import (
	"aula4/internal/repository/storage"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	MessageProductCreated = "Product created"
	MessageProductUpdated = "Product updated"
	MessageProductDeleted = "Product deleted"
)

type RequestBodyProduct struct {
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published *bool   `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

type Data struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

type ResponseBodyProduct struct {
	Message string `json:"message"`
	Data    *Data  `json:"data,omitempty"`
	Error   bool   `json:"error"`
}

type ResponseBodyTotalPrice struct {
	Products   []*Data `json:"products,omitempty"`
	TotalPrice float64 `json:"total_price"`
}

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

func ResponseWithError(w http.ResponseWriter, err error, statusCode int) {
	body := &ResponseBodyProduct{
		Message: http.StatusText(statusCode) + " - " + err.Error(),
		Data:    nil,
		Error:   true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func RespondWithProduct(w http.ResponseWriter, product *storage.Product, statusCode int, message string) {
	var body *ResponseBodyProduct
	if product == nil {
		body = &ResponseBodyProduct{
			Message: message,
			Data:    nil,
			Error:   false,
		}
	} else {
		dt := Data{
			Id:           product.Id,
			Name:         product.Name,
			Code_value:   product.Code_value,
			Is_published: *product.Is_published,
			Expiration:   product.Expiration,
			Quantity:     product.Quantity,
			Price:        product.Price,
		}

		body = &ResponseBodyProduct{
			Message: message,
			Data:    &dt,
			Error:   false,
		}
	}

	w.Header().Set("ßContent-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}
