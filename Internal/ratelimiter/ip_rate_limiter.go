package ratelimiter

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Kenasvarghese/Reverse-Proxy/internal/monitoring"
)

type rateLimit struct {
	lastUsed    atomic.Int64
	rateLimiter RateLimiter
}

type ipRateLimiter struct {
	rateLimiterConfig RateLimiterConfig
	ipRateLimit       sync.Map
	ob                *monitoring.Observer
	ttl               time.Duration
}

// NewIPRateLimiter creates a new IP rate limiter
func NewIPRateLimiter(rateLimiterConfig RateLimiterConfig, ob *monitoring.Observer) RateLimiter {
	rl := &ipRateLimiter{
		rateLimiterConfig: rateLimiterConfig,
		ipRateLimit:       sync.Map{},
		ob:                ob,
		ttl:               rateLimiterConfig.TTL,
	}
	go rl.CleanUp()
	return rl
}

// Allow checks if a request should be allowed based on available tokens.
func (p *ipRateLimiter) Allow(r *http.Request) bool {
	start := time.Now()
	ip := clientIP(r)
	now := time.Now().Unix()
	newRL := &rateLimit{
		lastUsed:    atomic.Int64{},
		rateLimiter: NewRateLimiter(p.rateLimiterConfig),
	}
	newRL.lastUsed.Store(now)
	value, loaded := p.ipRateLimit.LoadOrStore(ip, newRL)
	rl := value.(*rateLimit)
	rl.lastUsed.Store(now)
	if !loaded {
		p.ob.ActiveIpCount.Inc()
	}
	if !rl.rateLimiter.Allow(r) {
		p.ob.IpRateLimitedCount.Inc()
		p.ob.TimeForRatelimiterCheck.WithLabelValues("denied").Observe(float64(time.Since(start).Nanoseconds()))
		return false
	}
	p.ob.TimeForRatelimiterCheck.WithLabelValues("allowed").Observe(float64(time.Since(start).Nanoseconds()))
	return true
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
				p.ob.ActiveIpCount.Dec()
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
