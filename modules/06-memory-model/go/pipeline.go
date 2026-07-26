package conc

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// region:pipeline:start

// ParallelSquares squares each input concurrently with bounded parallelism,
// cancelling all workers on the first error. errgroup.WithContext gives the
// first-error-cancels behavior and SetLimit bounds concurrency; each worker
// writes a DISTINCT index, so there is no shared-write race. This is the pattern
// that replaces hand-rolled goroutine + channel + error plumbing.
func ParallelSquares(ctx context.Context, in []int, limit int) ([]int, error) {
	out := make([]int, len(in))
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(limit)
	for i, v := range in {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			out[i] = v * v
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return out, nil
}

// region:pipeline:end
