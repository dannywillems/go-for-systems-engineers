#!/usr/bin/env bash
# Build a SwiftPM product quietly and run it, so only the program's stdout is
# emitted (SwiftPM writes build progress that would otherwise pollute captured
# output). Usage: run-swift.sh <swift-dir> [product]
set -euo pipefail
dir="${1:?swift package dir}"
product="${2:-demo}"
cd "$dir"
swift build -c debug >/dev/null 2>&1
exec ".build/debug/$product"
