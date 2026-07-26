#!/usr/bin/env bash
# The falsifiability gate: regenerate every PORTABLE captured block and fail if
# any is stale. Uses `make docs-check`, whose -check pass skips portable:false
# blocks (assembly, timings, sizes) that legitimately differ between the
# author's arm64 and a CI runner. Requires all five toolchains on PATH; OPAM
# selects the CI opam switch.
set -euo pipefail

make docs-check OPAM="opam exec --"
