package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Kenasvarghese/Reverse-Proxy/Internal/config"
	"github.com/Kenasvarghese/Reverse-Proxy/Internal/middlewares"
	"github.com/Kenasvarghese/Reverse-Proxy/Internal/proxy"
	"github.com/Kenasvarghese/Reverse-Proxy/Internal/ratelimiter"
)

func main() {
	cfg := config.LoadConfig()
	rl := ratelimiter.NewRateLimiter(cfg.RateLimiterConfig)
	ipRL := ratelimiter.NewIPRateLimiter(cfg.RateLimiterConfig)
	proxyHandler := proxy.NewProxy(cfg.TransportConfig, cfg.GetOriginURL())
	wrappedHandler := middlewares.WrapHandler(proxyHandler,
		middlewares.RequestLogger,
		middlewares.GetRateLimiterMiddleware(ipRL),
		middlewares.GetRateLimiterMiddleware(rl),
	)
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
