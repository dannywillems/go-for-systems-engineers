#!/usr/bin/env bash
# Regenerate modules/09-reflection/measured.txt: the runtime cost of Go's
# reflection-based encoding/json versus a hand-written marshaler producing the
# identical output. MANUAL; the committed file is the source of truth. ns/op is
# non-deterministic (not run by `make docs`); allocs/op is deterministic and is
# the falsifiable part.
set -uo pipefail

export GOTOOLCHAIN=auto
cd "$(dirname "$0")/.." || exit 1
m=modules/09-reflection
out="$m/measured.txt"

{
	echo "machine: Apple M4 Pro (14 cores), macOS arm64"
	echo "workload: marshal a 2-field struct to JSON, two ways producing the same"
	echo "bytes. ReflectMarshal = encoding/json (walks reflect.Type per call);"
	echo "ManualMarshal = a hand-written encoder. allocs/op is deterministic."
	echo
	echo "# Go: reflection vs hand-written marshaler (go test -bench -benchmem)"
	(cd "$m/go" && go test -bench='Marshal' -benchmem -run='^$' ./... 2>/dev/null |
		grep -E 'Benchmark(Reflect|Manual)Marshal')
} >"$out" 2>/dev/null

echo "wrote $out"
cat "$out"
