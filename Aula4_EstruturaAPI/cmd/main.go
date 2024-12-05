package main

import (
	handlers "aula4/cmd/handler"
	"aula4/internal/domain"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/mod/sumdb/storage"
)

func main() {
	db := make(map[int]*domain.Product)
	//FillProductList(&db)

	ct := handlers.NewControllerProducts(storage: db)

	rt := chi.NewRouter()

	rt.Route("/products", func(r chi.Router) {
		r.Post("/", ct.Create)
		r.Get("/", ct.GetAll)
		r.Get("/{id}", ct.GetById)
		r.Get("/search", ct.Search)
	})

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
