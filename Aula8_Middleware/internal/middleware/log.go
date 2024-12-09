package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	size int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := &responseWriter{ResponseWriter: w}

		startTime := time.Now()
		next.ServeHTTP(wr, r)
		duration := time.Since(startTime)

		log.Printf("Method: %s, URL: %s, Duration: %v, Bytes: %d",
			r.Method,
			r.URL.String(),
			duration,
			wr.size,
		)
	})
}
