package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	rt := chi.NewRouter()
	rt.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
