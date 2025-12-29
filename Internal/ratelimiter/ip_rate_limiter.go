package ratelimiter

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type rateLimit struct {
	lastUsed    atomic.Int64
	rateLimiter RateLimiter
}

type ipRateLimiter struct {
	rateLimiterConfig RateLimiterConfig
	ipRateLimit       sync.Map
	ttl               time.Duration
}

// NewIPRateLimiter creates a new IP rate limiter
func NewIPRateLimiter(rateLimiterConfig RateLimiterConfig) RateLimiter {
	rl := &ipRateLimiter{
		rateLimiterConfig: rateLimiterConfig,
		ipRateLimit:       sync.Map{},
		ttl:               rateLimiterConfig.TTL,
	}
	go rl.CleanUp()
	return rl
}

// Allow checks if a request should be allowed based on available tokens.
func (p *ipRateLimiter) Allow(r *http.Request) bool {
	ip := clientIP(r)
	now := time.Now().Unix()
	newRL := &rateLimit{
		lastUsed:    atomic.Int64{},
		rateLimiter: NewRateLimiter(p.rateLimiterConfig),
	}
	newRL.lastUsed.Store(now)
	value, _ := p.ipRateLimit.LoadOrStore(ip, newRL)
	rl := value.(*rateLimit)
	rl.lastUsed.Store(now)
	return rl.rateLimiter.Allow(r)
}

// CleanUp removes expired rate limiters
func (p *ipRateLimiter) CleanUp() {
	cleanupInterval := max(p.ttl/2, time.Minute)
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		p.ipRateLimit.Range(func(key, value any) bool {
			if value.(*rateLimit).lastUsed.Load() < time.Now().Add(-p.ttl).Unix() {
				p.ipRateLimit.Delete(key)
			}
			return true
		})
	}
}

// clientIP returns the client IP address from the request
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
