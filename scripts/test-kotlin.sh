#!/usr/bin/env bash
# Compile a Kotlin module's lib+test with kotlinc and run the test main, which
# uses `check(...)` and exits non-zero (AssertionError) on failure. This avoids
# a JUnit dependency in the self-contained kotlinc build. Usage:
# test-kotlin.sh <kotlin-dir>
set -euo pipefail

if [ -d /opt/homebrew/opt/openjdk/bin ]; then
	export PATH="/opt/homebrew/opt/openjdk/bin:$PATH"
	export JAVA_HOME="/opt/homebrew/opt/openjdk"
fi

dir="${1:?kotlin module dir}"
cd "$dir"
mkdir -p build
kotlinc lib test -include-runtime -d build/test.jar >/dev/null 2>&1
exec java -jar build/test.jar
