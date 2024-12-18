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

func boolPtr(b bool) *bool {
	return &b
}

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name         string
		sendToken    bool
		input        RequestBodyProduct
		expectedErr  error
		expectedCode int
	}{
		{
			name:      "Successful creation",
			sendToken: true,
			input: RequestBodyProduct{
				Name:         "Product A",
				Quantity:     5,
				Code_value:   "123yy",
				Is_published: boolPtr(true),
				Expiration:   "01/01/2025",
				Price:        10.0,
			},
			expectedErr:  nil,
			expectedCode: http.StatusCreated,
		},
		{
			name:      "Missing name",
			sendToken: true,
			input: RequestBodyProduct{
				Quantity:     5,
				Code_value:   "123yy",
				Is_published: boolPtr(true),
				Expiration:   "01/01/2025",
				Price:        10.0,
			},
			expectedErr:  errors.New("name is required"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:      "Duplicated code_value",
			sendToken: true,
			input: RequestBodyProduct{
				Name:         "Product B",
				Quantity:     5,
				Code_value:   "123yy",
				Is_published: boolPtr(true),
				Expiration:   "01/01/2025",
				Price:        10.0,
			},
			expectedErr:  errors.New("the code_value must be unique"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:      "Invalid expiration date",
			sendToken: true,
			input: RequestBodyProduct{
				Name:         "Product C",
				Quantity:     5,
				Code_value:   "123zz",
				Is_published: boolPtr(true),
				Expiration:   "invalid-date",
				Price:        10.0,
			},
			expectedErr:  errors.New("invalid date"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:      "Negative price",
			sendToken: true,
			input: RequestBodyProduct{
				Name:         "Product D",
				Quantity:     5,
				Code_value:   "123xx",
				Is_published: boolPtr(true),
				Expiration:   "01/01/2025",
				Price:        -5.0,
			},
			expectedErr:  errors.New("price must be non-negative"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:      "Unauthorized",
			sendToken: false,
			input: RequestBodyProduct{
				Name:         "Product D",
				Quantity:     5,
				Code_value:   "123xx",
				Is_published: boolPtr(true),
				Expiration:   "01/01/2025",
				Price:        -5.0,
			},
			expectedErr:  errors.New("Unauthorized"),
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewRepositoryProductsMock()
			productService := service.NewServiceProducts(&mockRepo)
			productHandler := NewHandlerProducts(&productService)
			os.Setenv("TOKEN", "1234")

			if tt.name == "Duplicated code_value" {
				mockRepo.Products["1"] = &storage.Product{
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df012",
					Name:         "Product A",
					Quantity:     5,
					Code_value:   "123yy",
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				}
			}

			jsonBody, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/products", strings.NewReader(string(jsonBody)))

			if tt.sendToken {
				req.Header.Set("Token", "1234")
			}

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

func TestUpdateProduct(t *testing.T) {
	tests := []struct {
		name         string
		sendToken    bool
		productID    string
		updates      RequestBodyProduct
		initialData  map[string]*storage.Product
		expectedErr  error
		expectedCode int
	}{
		{
			name:      "Successful update",
			sendToken: true,
			productID: "684963bb-7172-48ad-aecd-cdca3f0df012",
			updates: RequestBodyProduct{
				Name:         "Product AA",
				Quantity:     5,
				Code_value:   "123yy",
				Is_published: boolPtr(true),
				Expiration:   "01/01/2025",
				Price:        10.0,
			},
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df012": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df012",
					Code_value:   "123yy",
					Name:         "Product A",
					Quantity:     5,
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
			},
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name:      "Product not found",
			sendToken: true,
			productID: "684963bb-7172-48ad-aecd-cdca3f0df013",
			updates: RequestBodyProduct{
				Name:         "Product AA",
				Quantity:     5,
				Code_value:   "123yy",
				Is_published: boolPtr(true),
				Expiration:   "01/01/2025",
				Price:        10.0,
			},
			initialData:  nil,
			expectedErr:  nil,
			expectedCode: http.StatusCreated,
		},
		{
			name:      "Invalid expiration date format",
			sendToken: true,
			productID: "684963bb-7172-48ad-aecd-cdca3f0df012",
			updates: RequestBodyProduct{
				Name:         "Product AA",
				Quantity:     5,
				Code_value:   "123yy",
				Is_published: boolPtr(true),
				Expiration:   "invalid-date",
				Price:        10.0,
			},
			initialData: map[string]*storage.Product{
				"1": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df034",
					Code_value:   "123yy",
					Name:         "Product A",
					Quantity:     5,
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
			},
			expectedErr:  errors.New("invalid date"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:      "Duplicated code_value",
			sendToken: true,
			productID: "684963bb-7172-48ad-aecd-cdca3f0df012",
			updates: RequestBodyProduct{
				Name:         "Product AA",
				Quantity:     5,
				Code_value:   "123xx",
				Is_published: boolPtr(true),
				Expiration:   "01/01/2025",
				Price:        10.0,
			},
			initialData: map[string]*storage.Product{
				"1": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df013",
					Code_value:   "123yy",
					Name:         "Product A",
					Quantity:     5,
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
				"2": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df033",
					Code_value:   "123xx",
					Name:         "Product B",
					Quantity:     3,
					Is_published: boolPtr(false),
					Expiration:   "01/01/2026",
					Price:        15.0,
				},
			},
			expectedErr:  errors.New("the code_value must be unique"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:      "Missing code_value",
			sendToken: true,
			productID: "684963bb-7172-48ad-aecd-cdca3f0df012",
			updates: RequestBodyProduct{
				Name:         "Product AA",
				Quantity:     5,
				Is_published: boolPtr(true),
				Expiration:   "01/01/2025",
				Price:        10.0,
			},
			initialData: map[string]*storage.Product{
				"1": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df019",
					Code_value:   "123yy",
					Name:         "Product A",
					Quantity:     5,
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
			},
			expectedErr:  errors.New("code_value is required"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:      "Unauthorized",
			sendToken: false,
			productID: "684963bb-7172-48ad-aecd-cdca3f0df012",
			updates: RequestBodyProduct{
				Name:         "Product AA",
				Quantity:     5,
				Code_value:   "123yy",
				Is_published: boolPtr(true),
				Expiration:   "01/01/2025",
				Price:        10.0,
			},
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df012": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df012",
					Code_value:   "123yy",
					Name:         "Product A",
					Quantity:     5,
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
			},
			expectedErr:  errors.New("Unauthorized"),
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewRepositoryProductsMock()
			for id, product := range tt.initialData {
				mockRepo.Products[id] = product
			}

			productService := service.NewServiceProducts(&mockRepo)
			productHandler := NewHandlerProducts(&productService)
			os.Setenv("TOKEN", "1234")

			jsonBody, _ := json.Marshal(tt.updates)
			req, _ := http.NewRequest("PUT", "/products/"+tt.productID, strings.NewReader(string(jsonBody)))

			if tt.sendToken {
				req.Header.Set("Token", "1234")
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(productHandler.UpdateOrCreate)
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

				if response.Data.Name != tt.updates.Name {
					t.Errorf("handler returned unexpected name: got %v want %v", response.Data.Name, tt.updates.Name)
				}
			}
		})
	}
}

func TestGetById(t *testing.T) {
	tests := []struct {
		name         string
		sendToken    bool
		productID    string
		initialData  map[string]*storage.Product
		expected     *storage.Product
		expectedErr  error
		expectedCode int
	}{
		{
			name:      "Successful retrieval",
			sendToken: true,
			productID: "684963bb-7172-48ad-aecd-cdca3f0df012",
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df012": {
					Id:         "684963bb-7172-48ad-aecd-cdca3f0df012",
					Name:       "Product A",
					Quantity:   5,
					Code_value: "123yy",
					Expiration: "01/01/2025",
					Price:      10.0,
				},
			},
			expected: &storage.Product{
				Id:         "684963bb-7172-48ad-aecd-cdca3f0df012",
				Name:       "Product A",
				Quantity:   5,
				Code_value: "123yy",
				Expiration: "01/01/2025",
				Price:      10.0,
			},
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name:      "Product not found",
			sendToken: true,
			productID: "684963bb-1111-48ad-aecd-cdca3f0df012",
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df012": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df012",
					Name:         "Product A",
					Quantity:     5,
					Code_value:   "123yy",
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
			},
			expected:     nil,
			expectedErr:  errors.New("product not found"),
			expectedCode: http.StatusNotFound,
		},
		{
			name:      "Product not found",
			sendToken: true,
			productID: "111",
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df012": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df012",
					Name:         "Product A",
					Quantity:     5,
					Code_value:   "123yy",
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
			},
			expected:     nil,
			expectedErr:  errors.New("product not found"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:      "Unauthorized",
			sendToken: false,
			productID: "111",
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df012": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df012",
					Name:         "Product A",
					Quantity:     5,
					Code_value:   "123yy",
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
			},
			expected:     nil,
			expectedErr:  errors.New("Unauthorized"),
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewRepositoryProductsMock()
			for id, product := range tt.initialData {
				mockRepo.Products[id] = product
			}

			productService := service.NewServiceProducts(&mockRepo)
			productHandler := NewHandlerProducts(&productService)
			os.Setenv("TOKEN", "1234")

			req, _ := http.NewRequest("GET", "/products/"+tt.productID, nil)

			if tt.sendToken {
				req.Header.Set("Token", "1234")
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(productHandler.GetById)
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
				var responseData storage.Product
				json.NewDecoder(rr.Body).Decode(&responseData)

				if responseData != *tt.expected {
					t.Errorf("handler returned unexpected product: got %+v want %+v", responseData, *tt.expected)
				}
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name          string
		sendToken     bool
		initialData   map[string]*storage.Product
		expectedCount int
		expectedCode  int
	}{
		{
			name:      "Successful retrieval of all products",
			sendToken: true,
			initialData: map[string]*storage.Product{
				"1": {
					Id:         "1",
					Name:       "Product A",
					Quantity:   5,
					Code_value: "123yy",
					Expiration: "01/01/2025",
					Price:      10.0,
				},
				"2": {
					Id:         "2",
					Name:       "Product B",
					Quantity:   10,
					Code_value: "456yy",
					Expiration: "01/01/2026",
					Price:      20.0,
				},
			},
			expectedCount: 2,
			expectedCode:  http.StatusOK,
		},
		{
			name:          "No Products",
			sendToken:     true,
			initialData:   map[string]*storage.Product{},
			expectedCount: 0,
			expectedCode:  http.StatusInternalServerError,
		},
		{
			name:      "Unauthorized",
			sendToken: false,
			initialData: map[string]*storage.Product{
				"1": {
					Id:         "1",
					Name:       "Product A",
					Quantity:   5,
					Code_value: "123yy",
					Expiration: "01/01/2025",
					Price:      10.0,
				},
				"2": {
					Id:         "2",
					Name:       "Product B",
					Quantity:   10,
					Code_value: "456yy",
					Expiration: "01/01/2026",
					Price:      20.0,
				},
			},
			expectedCount: 0,
			expectedCode:  http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewRepositoryProductsMock()
			for id, product := range tt.initialData {
				mockRepo.Products[id] = product
			}

			productService := service.NewServiceProducts(&mockRepo)
			productHandler := NewHandlerProducts(&productService)
			os.Setenv("TOKEN", "1234")

			req, _ := http.NewRequest("GET", "/products", nil)

			if tt.sendToken {
				req.Header.Set("Token", "1234")
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(productHandler.GetAll)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			if tt.expectedCount == 0 {
				var response ResponseBodyProduct
				json.NewDecoder(rr.Body).Decode(&response)

				if response.Error == false {
					t.Errorf("expected error but got none")
				}

				if response.Data != nil {
					t.Errorf("expected no data but got %v", response.Data)
				}
			} else {
				var response []storage.Product
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("error decoding response: %v", err)
				}

				if len(response) != tt.expectedCount {
					t.Errorf("expected %v products, got %v", tt.expectedCount, len(response))
				}
			}
		})
	}
}

func TestPatchProduct(t *testing.T) {
	tests := []struct {
		name         string
		sendToken    bool
		productID    string
		patchData    map[string]interface{}
		initialData  map[string]*storage.Product
		expectedErr  error
		expectedCode int
		expectedName string
	}{
		{
			name:      "Successful patch",
			sendToken: true,
			productID: "684963bb-7172-48ad-aecd-cdca3f0df012",
			patchData: map[string]interface{}{
				"name": "Product AA",
			},
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df012": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df012",
					Name:         "Product A",
					Quantity:     5,
					Code_value:   "123yy",
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
			},
			expectedErr:  nil,
			expectedCode: http.StatusOK,
			expectedName: "Product AA",
		},
		{
			name:      "Product not found",
			sendToken: true,
			productID: "684963bb-1212-48ad-aecd-cdca3f0df012",
			patchData: map[string]interface{}{
				"name": "Product AA",
			},
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df012": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df012",
					Name:         "Product A",
					Quantity:     5,
					Code_value:   "123yy",
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
			},
			expectedErr:  errors.New("product not found"),
			expectedCode: http.StatusNotFound,
			expectedName: "",
		},
		{
			name:      "Unauthorized",
			sendToken: false,
			productID: "684963bb-1212-48ad-aecd-cdca3f0df012",
			patchData: map[string]interface{}{
				"name": "Product AA",
			},
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df012": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df012",
					Name:         "Product A",
					Quantity:     5,
					Code_value:   "123yy",
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
			},
			expectedErr:  errors.New("Unauthorized"),
			expectedCode: http.StatusUnauthorized,
			expectedName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewRepositoryProductsMock()
			for id, product := range tt.initialData {
				mockRepo.Products[id] = product
			}

			productService := service.NewServiceProducts(&mockRepo)
			productHandler := NewHandlerProducts(&productService)
			os.Setenv("TOKEN", "1234")

			jsonBody, _ := json.Marshal(tt.patchData)
			req, _ := http.NewRequest("PATCH", "/products/"+tt.productID, strings.NewReader(string(jsonBody)))

			if tt.sendToken {
				req.Header.Set("Token", "1234")
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(productHandler.Update)
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

				if response.Data.Name != tt.expectedName {
					t.Errorf("handler returned unexpected name: got %v want %v", response.Data.Name, tt.expectedName)
				}
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tests := []struct {
		name         string
		sendToken    bool
		productID    string
		initialData  map[string]*storage.Product
		expectedErr  error
		expectedCode int
	}{
		{
			name:      "Successful deletion",
			sendToken: true,
			productID: "684963bb-7172-48ad-aecd-cdca3f0df012",
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df012": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df012",
					Name:         "Product A",
					Quantity:     5,
					Code_value:   "123yy",
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
			},
			expectedErr:  nil,
			expectedCode: http.StatusNoContent,
		},
		{
			name:      "Product not found",
			sendToken: true,
			productID: "684963bb-2323-48ad-aecd-cdca3f0df012",
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df012": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df012",
					Name:         "Product A",
					Quantity:     5,
					Code_value:   "123yy",
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
			},
			expectedErr:  errors.New("product not found"),
			expectedCode: http.StatusNotFound,
		},
		{
			name:      "Unauthorized",
			sendToken: false,
			productID: "684963bb-2323-48ad-aecd-cdca3f0df012",
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df012": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df012",
					Name:         "Product A",
					Quantity:     5,
					Code_value:   "123yy",
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
			},
			expectedErr:  errors.New("Unauthorized"),
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewRepositoryProductsMock()
			for id, product := range tt.initialData {
				mockRepo.Products[id] = product
			}

			productService := service.NewServiceProducts(&mockRepo)
			productHandler := NewHandlerProducts(&productService)
			os.Setenv("TOKEN", "1234")

			req, _ := http.NewRequest("DELETE", "/products/"+tt.productID, nil)

			if tt.sendToken {
				req.Header.Set("Token", "1234")
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(productHandler.Delete)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			if tt.expectedCode == http.StatusNoContent {
				_, exists := mockRepo.Products[tt.productID]
				if exists {
					t.Errorf("Product %v should have been deleted", tt.productID)
				}
			}
		})
	}
}
