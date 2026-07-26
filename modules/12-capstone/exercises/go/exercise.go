// Package exercises: Module 12 reader tasks. RED until you implement the stubs.
// Run the test with -race to prove your cache has no data race.
package exercises

import "context"

// Sem is a context-aware counting semaphore capping concurrent holders at n.
type Sem struct {
	ch chan struct{}
}

// NewSem returns a semaphore with n slots.
func NewSem(n int) *Sem {
	// TODO(reader): back the semaphore with a buffered channel of capacity n.
	return &Sem{}
}

// Acquire blocks until a slot is free OR ctx is done. It returns ctx.Err() if
// the context is cancelled while waiting, and nil once a slot is held.
func (s *Sem) Acquire(ctx context.Context) error {
	// TODO(reader): select on sending to s.ch versus <-ctx.Done().
	return nil
}

// Release returns a slot. Must be called exactly once per successful Acquire.
func (s *Sem) Release() {
	// TODO(reader): receive from s.ch.
}

// BoundedFetch returns the cached value for key, or fetches it via fetch and
// caches it. It must:
//   - be safe under concurrent callers (guard the map),
//   - never let the map exceed capacity (evict an arbitrary entry on overflow),
//   - bound concurrent fetches using sem, honoring ctx (return its error if the
//     context is cancelled while waiting for a slot),
//   - not hold the map lock across the fetch call.
//
// The provided mu guards m; sem is shared across callers.
func BoundedFetch(
	ctx context.Context,
	mu Locker, m map[int]int, capacity int,
	sem *Sem, key int, fetch func(int) int,
) (int, error) {
	// TODO(reader): implement the get-or-fetch path described above.
	return 0, nil
}

// Locker is the subset of sync.Locker the exercise needs (so the test can pass
// a *sync.Mutex without importing sync here).
type Locker interface {
	Lock()
	Unlock()
}
