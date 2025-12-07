package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Kenasvarghese/Caching-Proxy/Internal/config"
	"github.com/Kenasvarghese/Caching-Proxy/Internal/middlewares"
	"github.com/Kenasvarghese/Caching-Proxy/Internal/proxy"
	"github.com/Kenasvarghese/Caching-Proxy/Internal/rate_limiter"
)

func main() {
	cfg := config.LoadConfig()
	rl := rate_limiter.NewRateLimiter(cfg.RateLimiterConfig)
	proxyHandler := proxy.NewProxy(cfg.TransportConfig, cfg.GetOriginURL())
	wrappedHandler := middlewares.WrapHandler(proxyHandler,
		middlewares.RequestLogger,
		middlewares.GetRateLimiterMiddleware(rl))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: wrappedHandler,
	}
	log.Printf("proxy is listening on port %d", cfg.Port)
	log.Printf("forward to origin is %s", cfg.Origin)
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Error listening on port %d: %v", cfg.Port, err)
	}
}
