#!/usr/bin/env bash
# Compile a Kotlin module's lib+app with kotlinc and run it (no Gradle), so only
# the program's stdout is emitted. Convention: lib/ = library (no main), app/ =
# the demo main. Usage: run-kotlin.sh <kotlin-dir>
set -euo pipefail

# Homebrew's openjdk is keg-only; add it if present (CI provides java on PATH).
if [ -d /opt/homebrew/opt/openjdk/bin ]; then
	export PATH="/opt/homebrew/opt/openjdk/bin:$PATH"
	export JAVA_HOME="/opt/homebrew/opt/openjdk"
fi

dir="${1:?kotlin module dir}"
cd "$dir"
mkdir -p build
kotlinc lib app -include-runtime -d build/demo.jar >/dev/null 2>&1
exec java -jar build/demo.jar
