package main

import (
	"aula4/internal/handler"
	"aula4/internal/middleware"
	"aula4/internal/repository"
	"aula4/internal/repository/storage"
	"aula4/internal/service"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	st := storage.NewStorageProducts()
	rp := repository.NewRepositoryProducts(&st)
	sv := service.NewServiceProducts(&rp)
	hd := handler.NewHandlerProducts(&sv)

	rt := chi.NewRouter()

	rt.Use(middleware.LoggingMiddleware)
	rt.Use(middleware.ValidateToken)

	rt.Route("/products", func(r chi.Router) {
		r.Get("/", hd.GetAll)
		r.Get("/{id}", hd.GetById)
		r.Get("/search", hd.Search)
		r.Get("/consumer_price", hd.ConsumerPrice)
		r.Post("/", hd.Create)
		r.Put("/{id}", hd.UpdateOrCreate)
		r.Patch("/{id}", hd.Update)
		r.Delete("/{id}", hd.Delete)
	})

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
