package handler

import (
	"aula4/internal/repository"
	"aula4/internal/repository/storage"
	"aula4/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func boolPtrrrrr(b bool) *bool {
	return &b
}

func TestCreateProduect(t *testing.T) {
	tests := []struct {
		name         string
		input        RequestBodyProduct
		expectedErr  error
		expectedCode int
	}{
		{
			name: "Successful creation",
			input: RequestBodyProduct{
				Name:         "Product A",
				Quantity:     5,
				Code_value:   "123yy",
				Is_published: boolPtrrrrr(true),
				Expiration:   "01/01/2025",
				Price:        10.0,
			},
			expectedErr:  nil,
			expectedCode: http.StatusCreated,
		},
		{
			name: "Missing name",
			input: RequestBodyProduct{
				Quantity:     5,
				Code_value:   "123yy",
				Is_published: boolPtrrrrr(true),
				Expiration:   "01/01/2025",
				Price:        10.0,
			},
			expectedErr:  errors.New("name is required"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Duplicated code_value",
			input: RequestBodyProduct{
				Name:         "Product B",
				Quantity:     5,
				Code_value:   "123yy",
				Is_published: boolPtrrrrr(true),
				Expiration:   "01/01/2025",
				Price:        10.0,
			},
			expectedErr:  errors.New("the code_value must be unique"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid expiration date",
			input: RequestBodyProduct{
				Name:         "Product C",
				Quantity:     5,
				Code_value:   "123zz",
				Is_published: boolPtrrrrr(true),
				Expiration:   "invalid-date",
				Price:        10.0,
			},
			expectedErr:  errors.New("invalid date"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Negative price",
			input: RequestBodyProduct{
				Name:         "Product D",
				Quantity:     5,
				Code_value:   "123xx",
				Is_published: boolPtrrrrr(true),
				Expiration:   "01/01/2025",
				Price:        -5.0,
			},
			expectedErr:  errors.New("price must be non-negative"),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewRepositoryProductsMock()
			productService := service.NewServiceProducts(&mockRepo)
			productHandler := NewHandlerProducts(&productService)
			os.Setenv("TOKEN", "1234")

			// Adiciona um produto para verificar a duplicação
			if tt.name == "Duplicated code_value" {
				mockRepo.Products["1"] = &storage.Product{
					Id:           "1",
					Name:         "Product A",
					Quantity:     5,
					Code_value:   "123yy",
					Is_published: boolPtrrrrr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				}
			}

			jsonBody, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/products", strings.NewReader(string(jsonBody)))
			req.Header.Set("Token", "1234")
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(productHandler.Create)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			if tt.expectedErr != nil {

				var response ResponseBodyProduct
				json.NewDecoder(rr.Body).Decode(&response)

				if response.Error == false {
					t.Errorf("expected error but got none")
				}

				if response.Data != nil {
					t.Errorf("expected no data but got %v", response.Data)
				}
			} else {
				var response ResponseBodyProduct
				json.NewDecoder(rr.Body).Decode(&response)

				if response.Data.Name != tt.input.Name {
					t.Errorf("handler returned unexpected body: got %v want %v", response.Data.Name, tt.input.Name)
				}
			}
		})
	}
}

func TestUpdatePreoduct(t *testing.T) {
	mockRepo := repository.NewRepositoryProductsMock()
	productService := service.NewServiceProducts(&mockRepo)
	productHandler := NewHandlerProducts(&productService)

	os.Setenv("TOKEN", "1234")

	productID := "1"
	newProduct := RequestBodyProduct{
		Name:         "Product A",
		Quantity:     5,
		Code_value:   "123yy",
		Is_published: boolPtrrrrr(true),
		Expiration:   "01/01/2025",
		Price:        10.0,
	}

	mockRepo.Products[productID] = &storage.Product{
		Id:           productID,
		Name:         newProduct.Name,
		Quantity:     newProduct.Quantity,
		Code_value:   newProduct.Code_value,
		Is_published: newProduct.Is_published,
		Expiration:   newProduct.Expiration,
		Price:        newProduct.Price,
	}

	updates := RequestBodyProduct{
		Name:         "Product AA",
		Quantity:     5,
		Code_value:   "123yy",
		Is_published: boolPtrrrrr(true),
		Expiration:   "01/01/2025",
		Price:        10.0,
	}

	jsonBody, _ := json.Marshal(updates)

	req, _ := http.NewRequest("PUT", "/products/"+productID, strings.NewReader(string(jsonBody)))
	req.Header.Set("Token", "1234")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(productHandler.UpdateOrCreate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response ResponseBodyProduct
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Data.Name != updates.Name {
		t.Errorf("handler returned unexpected name: got %v want %v", response.Data.Name, updates.Name)
	}

	if response.Data.Quantity != updates.Quantity {
		t.Errorf("handler returned unexpected quantity: got %v want %v", response.Data.Quantity, updates.Quantity)
	}

	if response.Data.Code_value != updates.Code_value {
		t.Errorf("handler returned unexpected code value: got %v want %v", response.Data.Code_value, updates.Code_value)
	}

	if response.Data.Is_published != *updates.Is_published {
		t.Errorf("handler returned unexpected publish status: got %v want %v", response.Data.Is_published, *updates.Is_published)
	}

	if response.Data.Expiration != updates.Expiration {
		t.Errorf("handler returned unexpected expiration: got %v want %v", response.Data.Expiration, updates.Expiration)
	}

	if response.Data.Price != updates.Price {
		t.Errorf("handler returned unexpected price: got %v want %v", response.Data.Price, updates.Price)
	}
}

func TestGetBeyId(t *testing.T) {
	mockRepo := repository.NewRepositoryProductsMock()
	productService := service.NewServiceProducts(&mockRepo)
	productHandler := NewHandlerProducts(&productService)

	os.Setenv("TOKEN", "1234")

	productID := "2"
	newProduct := RequestBodyProduct{
		Name:         "Product A",
		Quantity:     5,
		Code_value:   "123ydy",
		Is_published: boolPtrrrrr(true),
		Expiration:   "01/01/2025",
		Price:        10.0,
	}

	mockRepo.Products[productID] = &storage.Product{
		Id:           productID,
		Name:         newProduct.Name,
		Quantity:     newProduct.Quantity,
		Code_value:   newProduct.Code_value,
		Is_published: newProduct.Is_published,
		Expiration:   newProduct.Expiration,
		Price:        newProduct.Price,
	}

	req, _ := http.NewRequest("GET", "/products/"+productID, nil)
	req.Header.Set("Token", "1234")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(productHandler.GetById)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response storage.Product
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Name != newProduct.Name {
		t.Errorf("handler returned unexpected body: got %v want %v", response.Name, newProduct.Name)
	}
}

func TestGetAell(t *testing.T) {
	mockRepo := repository.NewRepositoryProductsMock()
	productService := service.NewServiceProducts(&mockRepo)
	productHandler := NewHandlerProducts(&productService)

	os.Setenv("TOKEN", "1234")

	mockRepo.Products["1"] = &storage.Product{
		Id:           "1",
		Name:         "Product A",
		Quantity:     5,
		Code_value:   "123yy",
		Is_published: boolPtrrrrr(true),
		Expiration:   "01/01/2025",
		Price:        10.0,
	}

	mockRepo.Products["2"] = &storage.Product{
		Id:           "2",
		Name:         "Product B",
		Quantity:     10,
		Code_value:   "456yy",
		Is_published: boolPtrrrrr(false),
		Expiration:   "01/01/2026",
		Price:        20.0,
	}

	req, _ := http.NewRequest("GET", "/products", nil)
	req.Header.Set("Token", "1234")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(productHandler.GetAll)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response []storage.Product
	json.NewDecoder(rr.Body).Decode(&response)

	if len(response) != 2 {
		t.Errorf("expected 2 products, got %v", len(response))
	}
}

func TestPatchProdeuct(t *testing.T) {
	mockRepo := repository.NewRepositoryProductsMock()
	productService := service.NewServiceProducts(&mockRepo)
	productHandler := NewHandlerProducts(&productService)

	os.Setenv("TOKEN", "1234")

	productID := "1"
	mockRepo.Products[productID] = &storage.Product{
		Id:           productID,
		Name:         "Product A",
		Quantity:     5,
		Code_value:   "123yy",
		Is_published: boolPtrrrrr(true),
		Expiration:   "01/01/2025",
		Price:        10.0,
	}

	patchData := map[string]interface{}{
		"name": "Product AA",
	}

	jsonBody, _ := json.Marshal(patchData)

	req, _ := http.NewRequest("PATCH", "/products/"+productID, strings.NewReader(string(jsonBody)))
	req.Header.Set("Token", "1234")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(productHandler.Update)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response ResponseBodyProduct
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Data.Name != patchData["name"] {
		t.Errorf("handler returned unexpected name: got %v want %v", response.Data.Name, patchData["name"])
	}
}

func TestDeleteProdeuct(t *testing.T) {
	mockRepo := repository.NewRepositoryProductsMock()
	productService := service.NewServiceProducts(&mockRepo)
	productHandler := NewHandlerProducts(&productService)

	os.Setenv("TOKEN", "1234")

	productID := "1"
	mockRepo.Products[productID] = &storage.Product{
		Id:           productID,
		Name:         "Product A",
		Quantity:     5,
		Code_value:   "123yy",
		Is_published: boolPtrrrrr(true),
		Expiration:   "01/01/2025",
		Price:        10.0,
	}

	req, _ := http.NewRequest("DELETE", "/products/"+productID, nil)
	req.Header.Set("Token", "1234")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(productHandler.Delete)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	_, exists := mockRepo.Products[productID]
	if exists {
		t.Errorf("Product %v should have been deleted", productID)
	}
}
