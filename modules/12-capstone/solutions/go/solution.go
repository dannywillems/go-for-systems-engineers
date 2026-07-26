// Package solutions is the corrigé for Module 12. Run via `make solutions M=12`
// (and `go test -race`).
package solutions

import "context"

// Sem is a context-aware counting semaphore capping concurrent holders at n.
type Sem struct {
	ch chan struct{}
}

func NewSem(n int) *Sem { return &Sem{ch: make(chan struct{}, n)} }

func (s *Sem) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Sem) Release() { <-s.ch }

// Locker is the subset of sync.Locker the exercise needs.
type Locker interface {
	Lock()
	Unlock()
}

func BoundedFetch(
	ctx context.Context,
	mu Locker, m map[int]int, capacity int,
	sem *Sem, key int, fetch func(int) int,
) (int, error) {
	mu.Lock()
	if v, ok := m[key]; ok {
		mu.Unlock()
		return v, nil
	}
	mu.Unlock()

	if err := sem.Acquire(ctx); err != nil {
		return 0, err
	}
	defer sem.Release()

	v := fetch(key) // NOT under the map lock

	mu.Lock()
	if len(m) >= capacity {
		for k := range m { // evict one arbitrary entry
			delete(m, k)
			break
		}
	}
	m[key] = v
	mu.Unlock()
	return v, nil
}
