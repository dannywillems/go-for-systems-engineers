#!/usr/bin/env bash
# Run a Go benchmark set repeatedly and summarize with benchstat.
#
# Microbenchmarks lie. This wrapper defends against that as far as a shared
# machine allows: -count gives benchstat enough samples to report a confidence
# interval (the +/- column) instead of a single noisy number, and -run '^$'
# stops normal tests from perturbing timing. It does NOT make a laptop a
# quiet benchmarking host; treat the +/- and the "~" (not-significant) markers
# as the real signal, not the point estimate.
#
# Usage: go-bench.sh <go-module-dir> <bench-regex> <count> <summary-out>
set -euo pipefail

dir="${1:?go module dir}"
pattern="${2:-.}"
count="${3:-10}"
out="${4:?summary output path}"

raw="${out%.txt}.raw.txt"
export GOTOOLCHAIN=auto

( cd "$dir" && go test -run '^$' -bench "$pattern" -benchmem -count="$count" ) \
  | tee "$raw"

# Label the column "go" (benchstat otherwise prints the raw file path).
benchstat "go=$raw" > "$out"
echo "wrote $out"
