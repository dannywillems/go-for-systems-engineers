package testkit

import (
	"testing"
	"testing/synctest"
	"time"
)

// region:synctest:start

// TestLimiterRefill runs inside a synctest bubble: time is VIRTUAL, so the
// Sleep is instantaneous and the test is deterministic -- no real wall-clock
// wait, no flakiness. This is the Go 1.25 way to test time-dependent code.
func TestLimiterRefill(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := NewLimiter(2, 100*time.Millisecond)
		first, second := l.Allow(), l.Allow() // consume both tokens
		if !first || !second {
			t.Fatal("first two Allow should succeed")
		}
		if l.Allow() {
			t.Fatal("third Allow should fail (bucket empty)")
		}
		time.Sleep(100 * time.Millisecond) // virtual time: instant
		if !l.Allow() {
			t.Fatal("Allow should succeed after the window refills")
		}
	})
}

// region:synctest:end
