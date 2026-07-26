package observability

import (
	"io"
	"runtime/metrics"
	"runtime/pprof"
)

// region:pprof:start

// CPUProfile runs work while a CPU profile is being collected, writing the
// pprof-format profile to w. The profile is analyzed offline with
// `go tool pprof` (top, list, web) — profiling is a standard-library concern,
// not a third-party one.
func CPUProfile(w io.Writer, work func()) error {
	if err := pprof.StartCPUProfile(w); err != nil {
		return err
	}
	defer pprof.StopCPUProfile()
	work()
	return nil
}

// region:pprof:end

// region:metrics:start

// ReadMetric reads one sample from the stable runtime/metrics API. Unlike the
// legacy runtime.MemStats, this API is versioned and self-describing: every
// metric has a documented name like "/gc/heap/allocs:bytes" and a kind.
func ReadMetric(name string) uint64 {
	sample := []metrics.Sample{{Name: name}}
	metrics.Read(sample)
	if sample[0].Value.Kind() == metrics.KindUint64 {
		return sample[0].Value.Uint64()
	}
	return 0
}

// region:metrics:end
