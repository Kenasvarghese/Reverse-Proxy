package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Kenasvarghese/Reverse-Proxy/internal/config"
	"github.com/Kenasvarghese/Reverse-Proxy/internal/ebpf"
	"github.com/Kenasvarghese/Reverse-Proxy/internal/middlewares"
	"github.com/Kenasvarghese/Reverse-Proxy/internal/monitoring"
	"github.com/Kenasvarghese/Reverse-Proxy/internal/proxy"
	"github.com/Kenasvarghese/Reverse-Proxy/internal/ratelimiter"
	"github.com/cilium/ebpf/link"
)

func main() {
	cfg := config.LoadConfig()
	log.Println(cfg)
	ob := monitoring.NewObserver()
	beeObj := &ebpf.XDPObjects{}
	if err := ebpf.LoadXDPObjects(beeObj, nil); err != nil {
		panic(err)
	}
	defer beeObj.Close()
	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		log.Fatal(err)
	}
	link, err := link.AttachXDP(link.XDPOptions{
		Program:   beeObj.XdpBlockIp,
		Interface: iface.Index,
	})
	if err != nil {
		panic(err)
	}
	defer link.Close()

	rl := ratelimiter.NewRateLimiter(cfg.RateLimiterConfig)
	ipRL := ratelimiter.NewIPRateLimiter(cfg.RateLimiterConfig, beeObj.BlockedIps, ob)
	proxyHandler := proxy.NewProxy(cfg.TransportConfig, cfg.GetOriginURL())
	wrappedHandler := middlewares.WrapHandler(proxyHandler,
		middlewares.GetObservabilityMiddleware(ob),
		middlewares.RequestLogger,
		middlewares.GetRateLimiterMiddleware(ipRL),
		middlewares.GetRateLimiterMiddlewareWithObserver(rl, ob),
	)
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: wrappedHandler,
	}
	log.Printf("proxy is listening on port %d", cfg.Port)
	log.Printf("forward to origin is %s", cfg.Origin)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Error listening on port %d: %v", cfg.Port, err)
	}
}
