package handler

import (
	"aula4/internal/repository"
	"aula4/internal/service"
	"encoding/json"
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
	mockRepo := repository.NewRepositoryProductsMock()
	productService := service.NewServiceProducts(&mockRepo)
	productHandler := NewHandlerProducts(&productService)

	os.Setenv("TOKEN", "1234")

	newProduct := RequestBodyProduct{Name: "Product A", Quantity: 5, Code_value: "123yy", Is_published: boolPtr(true), Expiration: "01/01/2025", Price: 10.0}
	jsonBody, _ := json.Marshal(newProduct)

	req, _ := http.NewRequest("POST", "/products", strings.NewReader(string(jsonBody)))
	req.Header.Set("Token", "1234")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(productHandler.Create)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v"+rr.Body.String(), status, http.StatusCreated)
	}

	var response ResponseBodyProduct
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Data.Name != newProduct.Name {
		t.Errorf("handler returned unexpected body: got %v want %v", response.Data.Name, newProduct.Name)
	}
}

/*
func TestUpdateProduct(t *testing.T) {
	mockRepo := repository.NewRepositoryProductsMock()
	productService := service.NewServiceProducts(&mockRepo)
	productHandler := NewHandlerProducts(&productService)

	os.Setenv("TOKEN", "1234")

	// Criar o produto inicialmente
	newProduct := RequestBodyProduct{Name: "Product A", Quantity: 5, Code_value: "123", Is_published: boolPtr(true), Expiration: "2025-01-01", Price: 10.0}
	product := storage.Product{
		Id:           "1",
		Name:         newProduct.Name,
		Quantity:     newProduct.Quantity,
		Code_value:   newProduct.Code_value,
		Is_published: newProduct.Is_published,
		Expiration:   newProduct.Expiration,
		Price:        newProduct.Price,
	}

	mockRepo.Products[product.Id] = &product // Adiciona o produto ao mock

	// Agora tenta atualizar
	updates := map[string]interface{}{
		"name":     "Updated Product A",
		"quantity": 10,
		"price":    15.0,
	}
	jsonBody, _ := json.Marshal(updates)

	req, _ := http.NewRequest("PUT", "/products/1", strings.NewReader(string(jsonBody)))
	req.Header.Set("Token", "1234")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(productHandler.Update)
	handler.ServeHTTP(rr, req)

	// Verifique o status da resposta
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verifica a resposta
	var response ResponseBodyProduct
	json.NewDecoder(rr.Body).Decode(&response)

	// Verifica o nome do produto atualizado
	if response.Data.Name != "Updated Product A" {
		t.Errorf("handler returned unexpected body: got %v want %v", response.Data.Name, "Updated Product A")
	}

	// Verifica a quantidade do produto atualizado
	if response.Data.Quantity != 10 {
		t.Errorf("handler returned unexpected quantity: got %v want %v", response.Data.Quantity, 10)
	}
}
*/
