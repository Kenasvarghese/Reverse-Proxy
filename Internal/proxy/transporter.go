package proxy

import (
	"net"
	"net/http"
	"time"
)

// newTransporter creates an HTTP transport configured for high-concurrency proxying
// with aggressive connection pooling (100 conns/host) and strict timeouts
func newTransporter() http.RoundTripper {
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	customTransport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 2 * time.Second,
	}
	return customTransport
}
