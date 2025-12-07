package rate_limiter

import (
	"net/http"
	"sync"
	"time"
)

// tokenBucket implements the RateLimiter with token bucket algorithm for rate limiting.
type tokenBucket struct {
	bucketSize    uint64
	ratePerSecond float64
	bucket        uint64
	lastFilledAt  time.Time
	mu            sync.Mutex
}

// newTokenBucket creates a new token bucket rate limiter.
// bucketSize defines the maximum burst capacity.
// ratePerSecond defines how many requests are allowed per second.
func newTokenBucket(bucketSize uint64, ratePerSecond float64) RateLimiter {
	rl := &tokenBucket{
		bucketSize:    bucketSize,
		ratePerSecond: ratePerSecond,
		bucket:        bucketSize,
		lastFilledAt:  time.Now(),
		mu:            sync.Mutex{},
	}
	if ratePerSecond > 0 {
		rl.ratePerSecond = ratePerSecond
	}
	if bucketSize > 0 {
		rl.bucket = bucketSize
	}
	return rl
}

// Allow checks if a request should be allowed based on available tokens.
func (t *tokenBucket) Allow(r *http.Request) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.refill()
	if t.bucket > 0 {
		t.bucket--
		return true
	}
	return false
}

// refill adds tokens to the bucket based on elapsed time since last refill.
// Must be called with mutex held.
func (t *tokenBucket) refill() {
	timePassed := time.Since(t.lastFilledAt)
	tokensToAdd := uint64(timePassed.Seconds() * t.ratePerSecond)
	if tokensToAdd < 1 {
		return
	}
	t.bucket = min(tokensToAdd+t.bucket, t.bucketSize)
	t.lastFilledAt = time.Now()
}
