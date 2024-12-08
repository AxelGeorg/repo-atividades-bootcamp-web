package main

import (
	"aula4/internal/handler"
	"aula4/internal/repository"
	"aula4/internal/repository/storage"
	"aula4/internal/service"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	stor := storage.NewStorageProducts()
	repo := repository.NewRepositoryProducts(&stor)
	serv := service.NewServiceProducts(&repo)
	hand := handler.NewHandlerProducts(&serv)

	rt := chi.NewRouter()

	rt.Route("/products", func(r chi.Router) {
		r.Get("/", hand.GetAll)
		r.Get("/{id}", hand.GetById)
		r.Get("/search", hand.Search)
		r.Get("/consumer_price", hand.ConsumerPrice)
		r.Post("/", hand.Create)
		r.Put("/{id}", hand.UpdateOrCreate)
		r.Patch("/{id}", hand.Update)
		r.Delete("/{id}", hand.Delete)
	})

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
