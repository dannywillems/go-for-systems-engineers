// Package solutions is the corrigé for Module 10. Run via `make solutions M=10`.
package solutions

import (
	"sort"
	"strings"
)

func JoinOnce(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	total := len(sep) * (len(parts) - 1)
	for _, p := range parts {
		total += len(p)
	}
	var b strings.Builder
	b.Grow(total) // the single allocation
	for i, p := range parts {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(p)
	}
	return b.String()
}

func Percentile(xs []int64, p float64) int64 {
	if len(xs) == 0 {
		return 0
	}
	cp := make([]int64, len(xs))
	copy(cp, xs)
	sort.Slice(cp, func(i, j int) bool { return cp[i] < cp[j] })
	rank := int(p/100*float64(len(cp))+0.9999) - 1
	if rank < 0 {
		rank = 0
	}
	if rank >= len(cp) {
		rank = len(cp) - 1
	}
	return cp[rank]
}
