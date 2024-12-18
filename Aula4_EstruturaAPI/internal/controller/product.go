package controller

import (
	"aula4/internal/service/model"
	"aula4/internal/service/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type RequestBodyProduct struct {
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

type Data struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

type ProductController struct {
	ServiceProducts *model.ServiceProducts
}

func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBodyProduct
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		code := http.StatusBadRequest
		body := &ResponseBodyProduct{
			Message: "Bad Request",
			Data:    nil,
			Error:   true,
		}

		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}

	product := storage.Product{
		Name:         reqBody.Name,
		Quantity:     reqBody.Quantity,
		Code_value:   reqBody.Code_value,
		Is_published: reqBody.Is_published,
		Expiration:   reqBody.Expiration,
		Price:        reqBody.Price,
	}

	productServ, err := c.ServiceProducts.Create(product)
	if err != nil {
		code := http.StatusBadRequest
		body := &ResponseBodyProduct{
			Message: "Bad Request",
			Data:    nil,
			Error:   true,
		}

		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
		return
	}

	dt := Data{
		Id:           productServ.Id,
		Name:         productServ.Name,
		Code_value:   productServ.Code_value,
		Is_published: productServ.Is_published,
		Expiration:   productServ.Expiration,
		Quantity:     productServ.Quantity,
		Price:        productServ.Price,
	}

	code := http.StatusCreated
	body := &ResponseBodyProduct{
		Message: "Product created",
		Data:    &dt,
		Error:   false,
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func (c *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var products []*storage.Product
	for _, product := range c.ServiceProducts.Storage.DB {
		products = append(products, product)
	}
	json.NewEncoder(w).Encode(products)
}

func (c *ProductController) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	product := c.ServiceProducts.Storage.DB[idStr]
	if product != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
		return
	}

	http.NotFound(w, r)
}

func (c *ProductController) Search(w http.ResponseWriter, r *http.Request) {
	priceStr := r.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	var filteredProducts []*storage.Product
	for _, product := range c.ServiceProducts.Storage.DB {
		if product.Price > price {
			filteredProducts = append(filteredProducts, product)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredProducts)
}

func NewControllerProducts(service *model.ServiceProducts) *ProductController {
	return &ProductController{
		ServiceProducts: service,
	}
}
