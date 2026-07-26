#!/usr/bin/env bash
# Regenerate modules/12-capstone/measured.txt: the concurrent-cache throughput /
# latency comparison plus the migration matrix (binary size, core LOC, cold
# compile time). MANUAL; the committed file is the source of truth. Timings are
# non-deterministic and cold-compile times are especially noisy on a shared
# machine -- read them as order-of-magnitude, not precise.
set -uo pipefail

if [ -d /opt/homebrew/opt/openjdk/bin ]; then
	export PATH="/opt/homebrew/opt/openjdk/bin:$PATH"
fi
export GOTOOLCHAIN=auto

cd "$(dirname "$0")/.." || exit 1
m=modules/12-capstone
out="$m/measured.txt"

opam_exec() { if command -v opam >/dev/null; then opam exec -- "$@"; else "$@"; fi; }
# time a command in fractional seconds via python3 (BSD date lacks %N).
timeit() { python3 -c 'import subprocess,sys,time; t=time.time(); subprocess.run(sys.argv[1:],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL); print(f"{time.time()-t:.1f}")' "$@"; }

{
	echo "machine: Apple M4 Pro (14 cores), macOS arm64"
	echo "workload: 640k concurrent Get / 64 workers over a 256-entry bounded"
	echo "cache; backend fetch sleeps 100us; hot set == capacity so the cache"
	echo "absorbs ~99.9% of load. HTTP transport omitted (dependency-free)."
	echo
	echo "# Throughput and latency (single run, optimized)"
	(cd "$m/go" && go run ./cmd/bench)
	(cd "$m/rust" && cargo run --release --quiet --bin bench)
	(cd "$m/ocaml" && opam_exec dune exec --profile release bin/main.exe)
	(cd "$m/swift" && swift build -c release >/dev/null 2>&1 && .build/release/bench)
	bash scripts/run-kotlin.sh "$m/kotlin"
	echo
	echo "# Migration matrix"
	printf '%-8s %-14s %-10s %s\n' lang binary core-LOC cold-compile-s

	(cd "$m/go" && go build -o /tmp/cap_go ./cmd/bench)
	printf '%-8s %-14s %-10s %s\n' Go "$(wc -c </tmp/cap_go | tr -d ' ')" \
		"$(grep -cv '^[[:space:]]*$' "$m/go/cache.go")" \
		"$( (cd "$m/go" && go clean -cache >/dev/null 2>&1; timeit go build ./cmd/bench) )"

	(cd "$m/rust" && cargo build --release --quiet --bin bench >/dev/null 2>&1)
	printf '%-8s %-14s %-10s %s\n' Rust "$(wc -c <"$m/rust/target/release/bench" | tr -d ' ')" \
		"$(grep -cv '^[[:space:]]*$' "$m/rust/src/lib.rs")" \
		"$( (cd "$m/rust" && cargo clean >/dev/null 2>&1; timeit cargo build --release --bin bench) )"

	(cd "$m/ocaml" && opam_exec dune build --profile release >/dev/null 2>&1)
	printf '%-8s %-14s %-10s %s\n' OCaml \
		"$(wc -c <"$m/ocaml/_build/default/bin/main.exe" | tr -d ' ')" \
		"$(grep -cv '^[[:space:]]*$' "$m/ocaml/lib/capstone.ml")" \
		"$( (cd "$m/ocaml" && opam_exec dune clean >/dev/null 2>&1; timeit opam exec -- dune build --profile release) )"

	printf '%-8s %-14s %-10s %s\n' Swift \
		"$(wc -c <"$m/swift/.build/release/bench" | tr -d ' ')" \
		"$(grep -cv '^[[:space:]]*$' "$m/swift/Sources/Cache/Cache.swift")" \
		"$( (cd "$m/swift" && rm -rf .build >/dev/null 2>&1; timeit swift build -c release) )"

	printf '%-8s %-14s %-10s %s\n' Kotlin "$(wc -c <"$m/kotlin/build/demo.jar" | tr -d ' ')" \
		"$(grep -cv '^[[:space:]]*$' "$m/kotlin/lib/Cache.kt")" \
		"$(timeit kotlinc "$m/kotlin/lib" "$m/kotlin/app" -include-runtime -d /tmp/cap_k.jar)"
} >"$out" 2>/dev/null

echo "wrote $out"
cat "$out"
