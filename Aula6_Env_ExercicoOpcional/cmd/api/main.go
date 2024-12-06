package main

import (
	"aula4/internal/handler"
	"aula4/internal/repository"
	"aula4/internal/service/model"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	db := repository.NewMeliDB()
	serv := model.NewServiceProducts(db)
	ctrl := handler.NewControllerProducts(serv)

	rt := chi.NewRouter()

	rt.Route("/products", func(r chi.Router) {
		r.Get("/", ctrl.GetAll)
		r.Get("/{id}", ctrl.GetById)
		r.Get("/search", ctrl.Search)
		r.Get("/consumer_price", ctrl.ConsumerPrice)
		r.Post("/", ctrl.Create)
		r.Put("/{id}", ctrl.UpdateOrCreate)
		r.Patch("/{id}", ctrl.Update)
		r.Delete("/{id}", ctrl.Delete)
	})

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
