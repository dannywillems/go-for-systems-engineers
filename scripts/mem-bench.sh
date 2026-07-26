#!/usr/bin/env bash
# Regenerate modules/05-memory/measured.txt: the Go GOGC throughput/RSS trade
# and the cross-language allocation timings. MANUAL, like the benchmarks — the
# committed file is the source of truth, refreshed deliberately on a quiet
# machine (timings are non-deterministic, so this is not run by `make docs`).
set -uo pipefail

# Homebrew OpenJDK is keg-only; add it if present (for Kotlin/Java).
if [ -d /opt/homebrew/opt/openjdk/bin ]; then
	export PATH="/opt/homebrew/opt/openjdk/bin:$PATH"
fi
export GOTOOLCHAIN=auto

cd "$(dirname "$0")/.." || exit 1
m=modules/05-memory
out="$m/measured.txt"

opam_exec() { if command -v opam >/dev/null; then opam exec -- "$@"; else "$@"; fi; }

{
	echo "machine: Apple M4 Pro (14 cores), macOS arm64"
	echo "workload: 50,000,000 short-lived ~8-word heap objects, single run"
	echo
	echo "# Go: GOGC throughput vs RSS trade (same workload, two GOGC values)"
	echo "GOGC=100  $(cd "$m/go" && GOGC=100 go run ./cmd/gcdemo)"
	echo "GOGC=800  $(cd "$m/go" && GOGC=800 go run ./cmd/gcdemo)"
	echo
	echo "# Cross-language allocation time (optimized builds; see caveats)"
	(cd "$m/go" && GOGC=100 go run ./cmd/gcdemo | sed 's/^/Go     /')
	(cd "$m/ocaml" && opam_exec dune exec --profile release bin/main.exe)
	(cd "$m/rust" && cargo run --release --quiet --bin demo)
	(cd "$m/swift" && swift build -c release >/dev/null 2>&1 && .build/release/demo)
	bash scripts/run-kotlin.sh "$m/kotlin"
} >"$out" 2>/dev/null

echo "wrote $out"
cat "$out"
