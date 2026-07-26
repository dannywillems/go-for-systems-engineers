#!/usr/bin/env bash
# Fail if `go generate` output is not committed (Module 09 wires generators).
set -euo pipefail

while IFS= read -r modfile; do
	(cd "$(dirname "$modfile")" && go generate ./...)
done < <(find . -name go.mod -not -path '*/exercises/*' -not -path '*/solutions/*')

git diff --exit-code
