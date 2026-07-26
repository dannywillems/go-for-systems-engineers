package testkit

import (
	"sync"
	"time"
)

// Limiter is a token bucket: up to `burst` tokens, refilled `rate` per window.
// It uses time.Now directly, so a test can control its clock with
// testing/synctest (virtual time) rather than sleeping for real.
type Limiter struct {
	lastFill time.Time
	mu       sync.Mutex
	window   time.Duration
	tokens   int
	burst    int
}

func NewLimiter(burst int, window time.Duration) *Limiter {
	return &Limiter{tokens: burst, burst: burst, window: window, lastFill: time.Now()}
}

// Allow consumes a token if one is available, refilling the whole bucket once a
// full window has elapsed.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if time.Since(l.lastFill) >= l.window {
		l.tokens = l.burst
		l.lastFill = time.Now()
	}
	if l.tokens > 0 {
		l.tokens--
		return true
	}
	return false
}
