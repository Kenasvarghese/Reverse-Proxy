package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	proxy "github.com/Kenasvarghese/Caching-Proxy/proxy"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	origin := flag.String("origin", "", "URL to forward requests to")
	flag.Parse()
	if *origin == "" {
		log.Printf("forward to origin is empty")
		return
	}
	originURL, err := url.Parse(*origin)
	if err != nil {
		log.Printf("Error parsing URL: %v", err)
		return
	}
	proxyHandler := proxy.NewProxyHandler(originURL)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: proxyHandler,
	}
	log.Printf("proxy is listening on port %d", *port)
	log.Printf("forward to origin is %s", *origin)
	server.ListenAndServe()
}
