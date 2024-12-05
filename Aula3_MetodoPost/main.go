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
