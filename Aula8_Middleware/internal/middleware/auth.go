package middleware

import (
	"aula4/internal/utils"
	"errors"
	"net/http"
	"os"
)

func ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")
		if token == "" {
			utils.ResponseWithError(w, errors.New("authorization header is missing"), http.StatusUnauthorized)
			return
		}

		if token != os.Getenv("TOKEN") {
			utils.ResponseWithError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
