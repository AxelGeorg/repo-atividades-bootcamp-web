package main

import (
	"encoding/json"
	handlers "exercicio3/Handlers"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func FillProductList(db *map[int]*handlers.Product) {
	file, err := os.Open("products.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()

	var products []handlers.Product
	if err := json.NewDecoder(file).Decode(&products); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	for _, product := range products {
		(*db)[product.Id] = &product
	}
}

func main() {
	db := make(map[int]*handlers.Product)
	FillProductList(&db)

	ct := handlers.NewControllerProducts(db)

	rt := chi.NewRouter()

	rt.Route("/products", func(r chi.Router) {
		r.Post("/", ct.Create())
		r.Get("/", ct.GetAll())
		r.Get("/{id}", ct.GetById())
		r.Get("/search", ct.Search())
	})

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
