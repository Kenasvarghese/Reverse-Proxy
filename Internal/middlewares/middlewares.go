package middlewares

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

// WrapHandler wraps the handler with the provided middlewares
func WrapHandler(h http.Handler, middlewares ...Middleware) http.Handler {
	handler := h
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// RequestLogger is a middleware that logs the details of the request along with time taken to complete the request
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("Method: %s, Path: %s, Duration: %s", r.Method, r.URL.Path, time.Now().Sub(start))
	})
}
