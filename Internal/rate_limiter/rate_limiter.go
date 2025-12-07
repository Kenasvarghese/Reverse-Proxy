package rate_limiter

import (
	"net/http"
)

// RateLimiterConfig holds configuration for rate limiting behavior
type RateLimiterConfig struct {
	MaxAllowedRequests uint64 `default:"100" split_words:"true"`
	RequestRate        uint64 `default:"100" split_words:"true"`
	RateLimiterType    string `split_words:"true"`
}

const TokenBucketRateLimiterType = "TokenBucket"

type RateLimiter interface {
	Allow(r *http.Request) bool
}

// NewRateLimiter creates a rate limiter based on the specified configuration.
// Returns a token bucket rate limiter by default.
func NewRateLimiter(cfg RateLimiterConfig) RateLimiter {
	switch cfg.RateLimiterType {
	case TokenBucketRateLimiterType:
		return newTokenBucket(cfg.MaxAllowedRequests, cfg.RequestRate)
	default:
		return newTokenBucket(cfg.MaxAllowedRequests, cfg.RequestRate)
	}
}
