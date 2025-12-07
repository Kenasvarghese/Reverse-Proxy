package proxy

import (
	"net"
	"net/http"
	"time"
)

// TransportConfig holds HTTP transport configuration for the reverse proxy
type TransportConfig struct {
	// Transport config - Timeouts
	DialTimeout           time.Duration `default:"5s" split_words:"true"`
	ResponseHeaderTimeout time.Duration `default:"5s" split_words:"true"`
	TLSHandshakeTimeout   time.Duration `default:"10s" split_words:"true"`

	// Transport config - Connection pooling
	MaxIdleConns        int `default:"100" split_words:"true"`
	MaxIdleConnsPerHost int `default:"100" split_words:"true"`
	MaxConnsPerHost     int `default:"100" split_words:"true"`
}

// newTransporter creates an HTTP transport configured for reverse proxying.
func newTransporter(cfg TransportConfig) http.RoundTripper {
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	if cfg.DialTimeout > 0 {
		dialer.Timeout = cfg.DialTimeout
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
		ResponseHeaderTimeout: 5 * time.Second,
	}
	if cfg.MaxIdleConns > 0 {
		customTransport.MaxIdleConns = cfg.MaxIdleConns
	}
	if cfg.MaxIdleConnsPerHost > 0 {
		customTransport.MaxIdleConnsPerHost = cfg.MaxIdleConnsPerHost
	}
	if cfg.MaxConnsPerHost > 0 {
		customTransport.MaxConnsPerHost = cfg.MaxConnsPerHost
	}
	if cfg.ResponseHeaderTimeout > time.Duration(0) {
		customTransport.ResponseHeaderTimeout = cfg.ResponseHeaderTimeout
	}
	if cfg.TLSHandshakeTimeout > time.Duration(0) {
		customTransport.TLSHandshakeTimeout = cfg.TLSHandshakeTimeout
	}
	return customTransport
}
