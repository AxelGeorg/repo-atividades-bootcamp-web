package handlers

import (
	hand "aula4/cmd/handler"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ControllerProducts struct {
	storage map[int]*hand.Product
}

type RequestBodyProduct struct {
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

type ResponseBodyProduct struct {
	Message string       `json:"message"`
	Data    *DataProduct `json:"data"`
	Error   bool         `json:"error"`
}

type DataProduct struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

func NewControllerProducts(storage map[int]*hand.Product) *ControllerProducts {
	return &ControllerProducts{
		storage: storage,
	}
}

func (c *ControllerProducts) Create(w http.ResponseWriter, r *http.Request) {
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

	pr := &hand.Product{
		Id:           len(c.storage) + 1, // mudar isso aquiiiiiiiiii
		Name:         reqBody.Name,
		Quantity:     reqBody.Quantity,
		Code_value:   reqBody.Code_value,
		Is_published: reqBody.Is_published,
		Expiration:   reqBody.Expiration,
		Price:        reqBody.Price,
	}

	c.storage[pr.Id] = pr

	code := http.StatusCreated
	body := &ResponseBodyProduct{
		Message: "Product created",
		Data: &DataProduct{
			Id:           pr.Id,
			Name:         pr.Name,
			Quantity:     pr.Quantity,
			Code_value:   pr.Code_value,
			Is_published: pr.Is_published,
			Expiration:   pr.Expiration,
			Price:        pr.Price,
		},
		Error: false,
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func (c *ControllerProducts) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var products []*hand.Product
	for _, product := range c.storage {
		products = append(products, product)
	}
	json.NewEncoder(w).Encode(products)
}

func (c *ControllerProducts) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format"+err.Error(), http.StatusBadRequest)
		return
	}

	product := c.storage[id]
	if product != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
		return
	}

	http.NotFound(w, r)
}

func (c *ControllerProducts) Search(w http.ResponseWriter, r *http.Request) {
	priceStr := r.URL.Query().Get("price")
	fmt.Println("pegou:" + priceStr)
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	var filteredProducts []*hand.Product
	for _, product := range c.storage {
		if product.Price > price {
			filteredProducts = append(filteredProducts, product)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredProducts)
}
