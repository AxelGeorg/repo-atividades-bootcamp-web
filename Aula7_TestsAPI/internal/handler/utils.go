package handler

import (
	"aula4/internal/repository/storage"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func ValidateToken(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("Token")
	if token == "" {
		ResponseWithError(w, errors.New("authorization header is missing"), http.StatusUnauthorized)
		return false
	}

	if token != os.Getenv("TOKEN") {
		ResponseWithError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return false
	}

	return true
}
