#!/usr/bin/env bash
# Run every module's Go benchmarks (manual dispatch only). Results are
# INFORMATIONAL on shared CI runners; the committed bench.txt files, produced
# on a quiet machine via scripts/go-bench.sh, remain the source of truth.
set -uo pipefail

count="${1:-6}"
while IFS= read -r bench; do
	d="$(dirname "$bench")"
	echo "== $d"
	(cd "$d" && go test -run '^$' -bench . -benchmem -count="$count" || true)
done < <(find modules -name bench.txt)
