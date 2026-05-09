package src

import (
	"net/http"
	"sync"
	"time"
)

// ---------------- Rate Limiter ----------------
type RateLimiter struct {
	requests      int
	maxRequests   int
	window        time.Duration
	mu            sync.Mutex
	lastResetTime time.Time
}

func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests:   maxRequests,
		window:        window,
		lastResetTime: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.Sub(rl.lastResetTime) > rl.window {
		rl.requests = 0
		rl.lastResetTime = now
	}

	if rl.requests < rl.maxRequests {
		rl.requests++
		return true
	}

	return false
}

func RateLimitMiddleware(handler http.Handler, limiter *RateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("429 - Too Many Requests"))
			return
		}
		handler.ServeHTTP(w, r)
	})
}
