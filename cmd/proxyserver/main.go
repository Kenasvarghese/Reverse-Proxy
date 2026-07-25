package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Kenasvarghese/Reverse-Proxy/internal/config"
	"github.com/Kenasvarghese/Reverse-Proxy/internal/middlewares"
	"github.com/Kenasvarghese/Reverse-Proxy/internal/monitoring"
	"github.com/Kenasvarghese/Reverse-Proxy/internal/proxy"
	"github.com/Kenasvarghese/Reverse-Proxy/internal/ratelimiter"
)

func main() {
	cfg := config.LoadConfig()
	ob := monitoring.NewObserver()
	rl := ratelimiter.NewRateLimiter(cfg.RateLimiterConfig)
	ipRL := ratelimiter.NewIPRateLimiter(cfg.RateLimiterConfig, ob)
	proxyHandler := proxy.NewProxy(cfg.TransportConfig, cfg.GetOriginURL())
	wrappedHandler := middlewares.WrapHandler(proxyHandler,
		middlewares.GetObservabilityMiddleware(ob),
		middlewares.RequestLogger,
		middlewares.GetRateLimiterMiddleware(ipRL),
		middlewares.GetRateLimiterMiddlewareWithObserver(rl, ob),
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
