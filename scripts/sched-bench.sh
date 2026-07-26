#!/usr/bin/env bash
# Regenerate modules/07-scheduler/measured.txt: Go task-latency percentiles under
# CPU-bound load, and a cross-language CPU-bound parallel throughput comparison.
# MANUAL (timings are non-deterministic); the committed file is the source of
# truth, refreshed deliberately on a quiet machine.
set -uo pipefail

if [ -d /opt/homebrew/opt/openjdk/bin ]; then
	export PATH="/opt/homebrew/opt/openjdk/bin:$PATH"
fi
export GOTOOLCHAIN=auto

cd "$(dirname "$0")/.." || exit 1
m=modules/07-scheduler
out="$m/measured.txt"

opam_exec() { if command -v opam >/dev/null; then opam exec -- "$@"; else "$@"; fi; }

{
	echo "machine: Apple M4 Pro (14 cores), macOS arm64"
	echo
	echo "# Go scheduler: task-completion latency vs worker count (CPU-bound)"
	(cd "$m/go" && go run ./cmd/latency)
	echo
	echo "# Cross-language: CPU-bound parallel sqrt-sum, 400M terms (optimized)"
	(cd "$m/go" && go run ./cmd/throughput)
	(cd "$m/ocaml" && opam_exec dune exec --profile release bin/main.exe)
	(cd "$m/rust" && cargo run --release --quiet --bin demo)
	(cd "$m/swift" && swift build -c release >/dev/null 2>&1 && .build/release/demo)
	bash scripts/run-kotlin.sh "$m/kotlin"
} >"$out" 2>/dev/null

echo "wrote $out"
cat "$out"
