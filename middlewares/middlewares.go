package middlewares

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func WrapHandler(h http.Handler, middlewares ...Middleware) http.Handler {
	handler := h
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("Method: %s, Path: %s, Duration: %s", r.Method, r.URL.Path, time.Now().Sub(start))
	})
}
