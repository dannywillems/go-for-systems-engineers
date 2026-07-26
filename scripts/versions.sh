#!/usr/bin/env bash
# Print the exact version of every compiler, formatter, linter, and static
# analyzer this repo drives. Captured (as a non-portable block) so each
# COMPARISON.md records the toolchain that produced its numbers.
set -uo pipefail

line() { printf '%-22s %s\n' "$1" "$2"; }
ver() { "$@" 2>&1 | head -1; }

echo "# compilers"
# Builds pin `go 1.26.5` in every go.mod and fetch it via GOTOOLCHAIN, even
# though the installed driver may be older; report the toolchain that compiles.
line "go"           "$(GOTOOLCHAIN=go1.26.5 ver go version)"
line "rustc"        "$(ver rustc --version)"
line "cargo"        "$(ver cargo --version)"
line "ocaml"        "$(ver ocaml --version)"
line "swift"        "$(swift --version 2>&1 | grep -o 'Swift version [0-9.]*' | head -1)"
line "kotlinc"      "$(ver kotlinc -version)"
line "java"         "$(java -version 2>&1 | head -1)"

echo
echo "# formatters"
line "gofmt"        "bundled with go"
line "rustfmt"      "$(ver rustfmt --version)"
line "ocamlformat"  "$(ver ocamlformat --version)"
line "swift-format" "bundled with swift ($(ver swift format --version))"
line "ktlint"       "$(ver ktlint --version)"

echo
echo "# linters / static analysis"
line "go vet"       "bundled with go"
line "staticcheck"  "$(ver staticcheck -version)"
line "golangci-lint" "$(ver golangci-lint version)"
line "govulncheck"  "$(ver govulncheck -version | tail -1)"
line "clippy"       "$(ver cargo clippy --version)"
line "ocaml -w"     "compiler warnings-as-errors (no separate linter is standard)"
line "swift"        "compiler diagnostics + swift-format lint"
line "ktlint"       "lint + format for Kotlin"

echo
echo "# build / bench tooling"
line "dune"         "$(ver dune --version)"
line "benchstat"    "golang.org/x/perf (installed)"
line "hyperfine"    "$(ver hyperfine --version)"
