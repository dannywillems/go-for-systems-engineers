// Package exercises: Module 06 reader tasks. RED until you implement the stubs.
// Run the test with -race (`go test -race`) to prove your solution is not just
// correct but race-free.
package exercises

// TODO(reader): have `workers` goroutines each increment a shared counter `per`
// times, and return the total. It must be correct AND race-free (use
// sync/atomic or a sync.Mutex). Result must be workers*per.
func Increment(workers, per int) int {
	return 0 // replace me
}

// TODO(reader): fan-in merge. Return a channel that yields every value from all
// input channels, closing once all inputs are drained. Do not lose or duplicate
// values, and do not leak goroutines.
func Merge(chans ...<-chan int) <-chan int {
	out := make(chan int)
	close(out)
	return out // replace me
}
