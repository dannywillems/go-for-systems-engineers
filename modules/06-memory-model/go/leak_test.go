//go:build leakdemo

// Excluded from the normal suite (it deliberately leaks). The capture tool runs
// it with `-tags leakdemo` to show goleak catching a goroutine that never exits.
// In real code, goleak.VerifyTestMain guards every package against such leaks.
package conc

import (
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestGoroutineLeak(t *testing.T) {
	defer goleak.VerifyNone(t)
	// A goroutine with no exit path: the classic leak. It is not GC'd; it pins
	// everything it captured for the process lifetime.
	go func() {
		time.Sleep(time.Hour)
	}()
	time.Sleep(20 * time.Millisecond) // let it start so goleak sees it
}
