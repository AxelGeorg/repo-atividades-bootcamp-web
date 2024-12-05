package main

import (
	handlers "exercicio3/Handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	db := make(map[int]*handlers.Product)
	ct := handlers.NewControllerProducts(db)

	rt := chi.NewRouter()

	rt.Route("/products", func(rt chi.Router) {
		rt.Post("/", ct.Create())
		rt.Get("/", ct.GetAll())
		rt.Get("/{id}", ct.GetById())
		rt.Get("/search", ct.Search())
	})

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
