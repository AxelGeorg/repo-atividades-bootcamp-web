package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Product struct {
	Id           int
	Name         string
	Quantity     int
	Code_value   string
	Is_published bool
	Expiration   string
	Price        float64
}

type ControllerProducts struct {
	storage map[int]*Product
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
	Message string `json:"message"`
	Data    *struct {
		Id           int     `json:"id"`
		Name         string  `json:"name"`
		Quantity     int     `json:"quantity"`
		Code_value   string  `json:"code_value"`
		Is_published bool    `json:"is_published"`
		Expiration   string  `json:"expiration"`
		Price        float64 `json:"price"`
	} `json:"data"`
	Error bool `json:"error"`
}

func NewControllerProducts(storage map[int]*Product) *ControllerProducts {
	return &ControllerProducts{
		storage: storage,
	}
}

func (c *ControllerProducts) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		pr := &Product{
			Id:           len(c.storage) + 1,
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
			Data: &struct {
				Id           int     `json:"id"`
				Name         string  `json:"name"`
				Quantity     int     `json:"quantity"`
				Code_value   string  `json:"code_value"`
				Is_published bool    `json:"is_published"`
				Expiration   string  `json:"expiration"`
				Price        float64 `json:"price"`
			}{Id: pr.Id, Name: pr.Name, Code_value: pr.Code_value, Is_published: pr.Is_published, Expiration: pr.Expiration, Quantity: pr.Quantity, Price: pr.Price},
			Error: false,
		}

		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}
}

func (c *ControllerProducts) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var products []*Product
		for _, product := range c.storage {
			products = append(products, product)
		}
		json.NewEncoder(w).Encode(products)
	}
}

func (c *ControllerProducts) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		for _, product := range c.storage {
			if product.Id == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(product)
				return
			}
		}

		http.NotFound(w, r)
	}
}

func (c *ControllerProducts) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		priceStr := r.URL.Query().Get("price")
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			http.Error(w, "Invalid price format", http.StatusBadRequest)
			return
		}

		var filteredProducts []*Product
		for _, product := range c.storage {
			if product.Price > price {
				filteredProducts = append(filteredProducts, product)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(filteredProducts)
	}
}
