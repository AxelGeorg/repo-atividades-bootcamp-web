package handler

import (
	"aula4/internal/middleware"
	"aula4/internal/repository"
	"aula4/internal/repository/storage"
	"aula4/internal/service"
	"aula4/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool {
	return &b
}

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name         string
		input        utils.RequestBodyProduct
		expectedErr  error
		expectedCode int
	}{
		{
			name: "Successful creation",
			input: utils.RequestBodyProduct{
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
			name: "Missing name",
			input: utils.RequestBodyProduct{
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
			name: "Duplicated code_value",
			input: utils.RequestBodyProduct{
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
			name: "Invalid expiration date",
			input: utils.RequestBodyProduct{
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
			name: "Negative price",
			input: utils.RequestBodyProduct{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewRepositoryProductsMock()
			productService := service.NewServiceProducts(&mockRepo)
			productHandler := NewHandlerProducts(&productService)
			os.Setenv("TOKEN", "1234")

			if tt.name == "Duplicated code_value" {
				mockRepo.Products["684963bb-7172-48ad-aecd-cdca3f0df012"] = &storage.Product{
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
			req.Header.Set("Token", "1234")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(productHandler.Create)
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "handler returned wrong status code")

			var response utils.ResponseBodyProduct
			json.NewDecoder(rr.Body).Decode(&response)

			if tt.expectedErr != nil {
				require.True(t, response.Error, "expected error but got none")
				require.Nil(t, response.Data, "expected no data but got some")
			} else {
				require.Equal(t, tt.input.Name, response.Data.Name, "handler returned unexpected body")
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	tests := []struct {
		name         string
		productID    string
		updates      utils.RequestBodyProduct
		initialData  map[string]*storage.Product
		expectedErr  error
		expectedCode int
	}{
		{
			name:      "Successful update",
			productID: "684963bb-7172-48ad-aecd-cdca3f0df012",
			updates: utils.RequestBodyProduct{
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
			productID: "684963bb-7172-48ad-aecd-cdca3f0df013",
			updates: utils.RequestBodyProduct{
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
			productID: "684963bb-7172-48ad-aecd-cdca3f0df012",
			updates: utils.RequestBodyProduct{
				Name:         "Product AA",
				Quantity:     5,
				Code_value:   "123yy",
				Is_published: boolPtr(true),
				Expiration:   "invalid-date",
				Price:        10.0,
			},
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df034": {
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
			productID: "684963bb-7172-48ad-aecd-cdca3f0df012",
			updates: utils.RequestBodyProduct{
				Name:         "Product AA",
				Quantity:     5,
				Code_value:   "123xx",
				Is_published: boolPtr(true),
				Expiration:   "01/01/2025",
				Price:        10.0,
			},
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df013": {
					Id:           "684963bb-7172-48ad-aecd-cdca3f0df013",
					Code_value:   "123yy",
					Name:         "Product A",
					Quantity:     5,
					Is_published: boolPtr(true),
					Expiration:   "01/01/2025",
					Price:        10.0,
				},
				"684963bb-7172-48ad-aecd-cdca3f0df033": {
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
			productID: "684963bb-7172-48ad-aecd-cdca3f0df012",
			updates: utils.RequestBodyProduct{
				Name:         "Product AA",
				Quantity:     5,
				Is_published: boolPtr(true),
				Expiration:   "01/01/2025",
				Price:        10.0,
			},
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df019": {
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
			req.Header.Set("Token", "1234")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(productHandler.UpdateOrCreate)
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "handler returned wrong status code")

			var response utils.ResponseBodyProduct
			json.NewDecoder(rr.Body).Decode(&response)

			if tt.expectedErr != nil {
				require.True(t, response.Error, "expected error but got none")
				require.Nil(t, response.Data, "expected no data but got some")
			} else {
				require.Equal(t, tt.updates.Name, response.Data.Name, "handler returned unexpected name")
			}
		})
	}
}

func TestGetById(t *testing.T) {
	tests := []struct {
		name         string
		productID    string
		initialData  map[string]*storage.Product
		expected     *storage.Product
		expectedErr  error
		expectedCode int
	}{
		{
			name:      "Successful retrieval",
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
			req.Header.Set("Token", "1234")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(productHandler.GetById)
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "handler returned wrong status code")

			if tt.expectedErr != nil {
				var response utils.ResponseBodyProduct
				json.NewDecoder(rr.Body).Decode(&response)

				require.True(t, response.Error, "expected error but got none")
				require.Nil(t, response.Data, "expected no data but got some")
			} else {
				var responseData storage.Product
				err := json.NewDecoder(rr.Body).Decode(&responseData)
				require.NoError(t, err, "could not decode response body")

				require.Equal(t, *tt.expected, responseData, "handler returned unexpected product")
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name          string
		initialData   map[string]*storage.Product
		expectedCount int
		expectedCode  int
	}{
		{
			name: "Successful retrieval of all products",
			initialData: map[string]*storage.Product{
				"684963bb-7172-48ad-aecd-cdca3f0df019": {
					Id:         "684963bb-7172-48ad-aecd-cdca3f0df019",
					Name:       "Product A",
					Quantity:   5,
					Code_value: "123yy",
					Expiration: "01/01/2025",
					Price:      10.0,
				},
				"684963bb-1313-48ad-aecd-cdca3f0df019": {
					Id:         "684963bb-1313-48ad-aecd-cdca3f0df019",
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
			initialData:   map[string]*storage.Product{},
			expectedCount: 0,
			expectedCode:  http.StatusInternalServerError,
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
			req.Header.Set("Token", "1234")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(productHandler.GetAll)
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "handler returned wrong status code")

			if tt.expectedCount == 0 {
				var response utils.ResponseBodyProduct
				err := json.NewDecoder(rr.Body).Decode(&response)
				require.NoError(t, err, "could not decode response body")

				require.True(t, response.Error, "expected error but got none")
				require.Nil(t, response.Data, "expected no data but got some")
			} else {
				var response []storage.Product
				err := json.NewDecoder(rr.Body).Decode(&response)
				require.NoError(t, err, "error decoding response")

				require.Equal(t, tt.expectedCount, len(response), "expected different number of products")
			}
		})
	}
}

func TestPatchProduct(t *testing.T) {
	tests := []struct {
		name         string
		productID    string
		patchData    map[string]interface{}
		initialData  map[string]*storage.Product
		expectedErr  error
		expectedCode int
		expectedName string
	}{
		{
			name:      "Successful patch",
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
			req.Header.Set("Token", "1234")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(productHandler.Update)
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "handler returned wrong status code")

			if tt.expectedErr != nil {
				var response utils.ResponseBodyProduct
				err := json.NewDecoder(rr.Body).Decode(&response)
				require.NoError(t, err, "could not decode response body")

				require.True(t, response.Error, "expected error but got none")
				require.Nil(t, response.Data, "expected no data but got some")
			} else {
				var response utils.ResponseBodyProduct
				err := json.NewDecoder(rr.Body).Decode(&response)
				require.NoError(t, err, "could not decode response body")

				require.Equal(t, tt.expectedName, response.Data.Name, "handler returned unexpected name")
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tests := []struct {
		name         string
		productID    string
		initialData  map[string]*storage.Product
		expectedErr  error
		expectedCode int
	}{
		{
			name:      "Successful deletion",
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
			req.Header.Set("Token", "1234")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(productHandler.Delete)
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "handler returned wrong status code")

			if tt.expectedCode == http.StatusNoContent {
				_, exists := mockRepo.Products[tt.productID]
				require.False(t, exists, "Product %v should have been deleted", tt.productID)
			}
		})
	}
}

func TestMiddleware(t *testing.T) {
	mockRepo := repository.NewRepositoryProductsMock()
	productService := service.NewServiceProducts(&mockRepo)
	productHandler := NewHandlerProducts(&productService)
	handler := http.HandlerFunc(productHandler.Create)

	tests := []struct {
		name         string
		token        string
		expectedCode int
	}{
		{
			name:         "Missing Token",
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Invalid Token",
			token:        "invalid-token",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/products", nil)
			if tt.token != "" {
				req.Header.Set("Token", tt.token)
			}

			rr := httptest.NewRecorder()
			middleware.ValidateToken(handler).ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "handler returned wrong status code")
		})
	}
}
