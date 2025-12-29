package ratelimiter

import (
	"net/http"
	"time"
)

// RateLimiterConfig holds configuration for rate limiting behavior
type RateLimiterConfig struct {
	MaxAllowedRequests uint64        `default:"150" split_words:"true"`
	RequestRatePerSec  float64       `default:"50" split_words:"true"`
	RateLimiterType    string        `split_words:"true"`
	TTL                time.Duration `default:"1m" split_words:"true"`
}

type RateLimiter interface {
	Allow(r *http.Request) bool
}

const TokenBucketWithChannelRateLimiterType = "TokenBucketWithChannel"
const TokenBucketRateLimiterType = "TokenBucket"

// NewRateLimiter creates a rate limiter based on the specified configuration.
// Returns a token bucket rate limiter by default.
func NewRateLimiter(cfg RateLimiterConfig) RateLimiter {
	switch cfg.RateLimiterType {
	case TokenBucketRateLimiterType:
		return newTokenBucket(cfg.MaxAllowedRequests, cfg.RequestRatePerSec)
	case TokenBucketWithChannelRateLimiterType:
		return newTokenBucketWithChannel(cfg)
	default:
		return newTokenBucket(cfg.MaxAllowedRequests, cfg.RequestRatePerSec)
	}
}
