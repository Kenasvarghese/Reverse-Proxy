package ratelimiter

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Kenasvarghese/Reverse-Proxy/internal/monitoring"
	"github.com/cilium/ebpf"
)

type rateLimit struct {
	lastUsed    atomic.Int64
	rateLimiter RateLimiter
}

type ipRateLimiter struct {
	rateLimiterConfig RateLimiterConfig
	ipRateLimit       sync.Map
	beeIPMap          *ebpf.Map
	ob                *monitoring.Observer
	ttl               time.Duration
}

// NewIPRateLimiter creates a new IP rate limiter
func NewIPRateLimiter(rateLimiterConfig RateLimiterConfig, beeIPMap *ebpf.Map, ob *monitoring.Observer) RateLimiter {
	rl := &ipRateLimiter{
		rateLimiterConfig: rateLimiterConfig,
		ipRateLimit:       sync.Map{},
		beeIPMap:          beeIPMap,
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
	binaryIp := uint32(0)
	if len(ip) >= 4 {
		binaryIp = binary.BigEndian.Uint32(ip)
	}
	now := time.Now().Unix()
	newRL := &rateLimit{
		lastUsed:    atomic.Int64{},
		rateLimiter: NewRateLimiter(p.rateLimiterConfig),
	}
	newRL.lastUsed.Store(now)
	value, loaded := p.ipRateLimit.LoadOrStore(binaryIp, newRL)
	rl := value.(*rateLimit)
	rl.lastUsed.Store(now)
	if !loaded {
		fmt.Println(p.ipRateLimit)
		p.ob.ActiveIpCount.Inc()
	}
	if !rl.rateLimiter.Allow(r) {
		p.beeIPMap.Put(binaryIp, uint8(1))
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
				p.beeIPMap.Delete(key)
				p.ob.ActiveIpCount.Dec()
			}
			return true
		})
	}
}

// clientIP returns the client IP address from the request
func clientIP(r *http.Request) net.IP {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return net.ParseIP(strings.TrimSpace(strings.Split(xff, ",")[0])).To4()
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return net.ParseIP(r.RemoteAddr).To4()
	}
	return net.ParseIP(host).To4()

}
