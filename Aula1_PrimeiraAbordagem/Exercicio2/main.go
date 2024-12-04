package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GreetingRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//curl -X POST http://localhost:8080/greetings/ -d '{"firstName": "Axel", "lastName": "Georg"}' -H "Content-Type: application/json"

func main() {
	rt := chi.NewRouter()
	rt.Post("/greetings/", func(w http.ResponseWriter, r *http.Request) {
		var req GreetingRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		greeting := fmt.Sprintf("Hello %s %s", req.FirstName, req.LastName)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(greeting))
	})

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
