package main

import (
	"aula4/internal/controller"
	"aula4/internal/repository"
	"aula4/internal/service/model"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	db := repository.NewMeliDB()
	serv := model.NewServiceProducts(db)
	ctrl := controller.NewControllerProducts(serv)

	rt := chi.NewRouter()

	rt.Route("/products", func(r chi.Router) {
		r.Post("/", ctrl.Create)
		r.Get("/", ctrl.GetAll)
		r.Get("/{id}", ctrl.GetById)
		r.Get("/search", ctrl.Search)
	})

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
