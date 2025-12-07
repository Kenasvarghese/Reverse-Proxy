package rate_limiter

import (
	"context"
	"net/http"
	"time"
)

type tokenBucketWithChannel struct {
	bucket chan struct{}
}

// newTokenBucketWithChannel creates a new token bucket rate limiter with channel bucket.
// bucketSize defines the maximum burst capacity.
// ratePerSecond defines how many requests are allowed per second.
func newTokenBucketWithChannel(cfg RateLimiterConfig) RateLimiter {
	bucket := make(chan struct{}, cfg.MaxAllowedRequests)
	rl := &tokenBucketWithChannel{
		bucket: bucket,
	}
	//TODO: add cancel to graceful shutdown
	ctx, _ := context.WithCancel(context.Background())
	//filler goroutine empties the tokens concurrently
	go rl.filler(ctx, time.Duration(float64(time.Second)/cfg.RequestRatePerSec))
	return rl
}

// Allow checks if a request should be allowed based on available tokens.
func (t *tokenBucketWithChannel) Allow(r *http.Request) bool {
	select {
	//request is allowed if there is space in the bucket
	case t.bucket <- struct{}{}:
		return true
	default:
		return false
	}
}

// filler removes tokens from the channel at the given rate
func (t *tokenBucketWithChannel) filler(ctx context.Context, rate time.Duration) {
	ticker := time.NewTicker(rate)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		case <-t.bucket:
			continue
		default:
			continue
		}
	}
}
