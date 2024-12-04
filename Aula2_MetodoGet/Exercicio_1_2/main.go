package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Product struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

var Products []Product

func FillProductList() {
	file, err := os.Open("products.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()

	if err := json.NewDecoder(file).Decode(&Products); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
}

func GetPingPong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Products)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	for _, product := range Products {
		if product.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	http.NotFound(w, r) // Retorna 404 se o produto não for encontrado
}

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	priceStr := r.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	var filteredProducts []Product
	for _, product := range Products {
		if product.Price > price {
			filteredProducts = append(filteredProducts, product)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredProducts)
}

func main() {
	FillProductList()

	rt := chi.NewRouter()
	rt.Get("/ping", GetPingPong)
	rt.Get("/products", GetProducts)
	rt.Get("/products/{id}", GetProductByID)
	rt.Get("/products/search", SearchProducts)

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
