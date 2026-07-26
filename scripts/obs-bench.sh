#!/usr/bin/env bash
# Regenerate modules/10-observability/measured.txt: the cross-language
# allocation cost (deterministic counts + non-portable ns/op) and the Go CPU
# profile's hottest function. MANUAL; the committed file is the source of truth.
# Timings are non-deterministic, so this is not run by `make docs`.
set -uo pipefail

if [ -d /opt/homebrew/opt/openjdk/bin ]; then
	export PATH="/opt/homebrew/opt/openjdk/bin:$PATH"
fi
export GOTOOLCHAIN=auto

cd "$(dirname "$0")/.." || exit 1
m=modules/10-observability
out="$m/measured.txt"

opam_exec() { if command -v opam >/dev/null; then opam exec -- "$@"; else "$@"; fi; }

{
	echo "machine: Apple M4 Pro (14 cores), macOS arm64"
	echo "workload: build a 64-part string two ways (naive concat vs pre-sized),"
	echo "1,000,000 iterations, optimized builds."
	echo
	echo "# Allocation cost per build (counts are deterministic; ns/op is not)"
	(cd "$m/rust" && cargo run --release --quiet --bin allocs | tail -1)
	(cd "$m/ocaml" && opam_exec dune exec --profile release bin/main.exe | tail -1)
	bash scripts/run-kotlin.sh "$m/kotlin" | tail -1
	(cd "$m/swift" && swift build -c release >/dev/null 2>&1 && .build/release/bench)
	echo
	echo "# Go microbenchmark (go test -bench -benchmem): allocs/op deterministic"
	(cd "$m/go" && go test -bench='Concat|Builder' -benchmem -run='^$' ./... 2>/dev/null |
		grep -E 'Benchmark(ConcatPlus|BuilderGrow)')
	echo
	echo "# Go CPU profile: hottest function (go tool pprof -top)"
	(cd "$m/go" && go run ./cmd/profile >/dev/null 2>&1 &&
		go tool pprof -top -nodecount=1 cpu.prof 2>/dev/null |
		grep -E 'observability' | head -1
	rm -f "$m/go/cpu.prof" cpu.prof 2>/dev/null)
	echo
	echo "# Go runtime/metrics readout (names are a stable API; values are not)"
	(cd "$m/go" && go run ./cmd/metrics)
} >"$out" 2>/dev/null

echo "wrote $out"
cat "$out"
