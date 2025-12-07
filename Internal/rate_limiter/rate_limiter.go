package rate_limiter

import (
	"fmt"
	"net/http"
)

// RateLimiterConfig holds configuration for rate limiting behavior
type RateLimiterConfig struct {
	MaxAllowedRequests uint64 `default:"100" split_words:"true"`
	RequestRate        uint64 `default:"100" split_words:"true"`
	RateLimiterType    string `split_words:"true"`
}

type RateLimiter interface {
	Allow(r *http.Request) bool
}

// NewRateLimiter creates a rate limiter based on the specified configuration.
// Returns a token bucket rate limiter by default.
func NewRateLimiter(cfg RateLimiterConfig) RateLimiter {
	fmt.Println(cfg)
	switch cfg.RateLimiterType {
	case "Token Bucket":
		return newTokenBucket(cfg.MaxAllowedRequests, cfg.RequestRate)
	default:
		return newTokenBucket(cfg.MaxAllowedRequests, cfg.RequestRate)

	}
}
