package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/Kenasvarghese/Caching-Proxy/Internal/rate_limiter"
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
		// ResponseWriter wrapper for logging the written http status
		lw := newLoggingResponseWriter(w)
		next.ServeHTTP(lw, r)
		log.Printf("Method: %s, Path: %s, Duration: %s, Status: %d", r.Method, r.URL.Path, time.Since(start), lw.status)
	})
}

// GetRateLimiterMiddleware returns the middleware with the provided rate limiter
func GetRateLimiterMiddleware(rl rate_limiter.RateLimiter) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rl.Allow(r) {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			}
		})
	}
}
