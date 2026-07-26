#!/usr/bin/env bash
# The falsifiability gate: regenerate every captured README block from real
# programs, then fail if anything drifted from what is committed. Uses the CI
# opam switch via the OPAM override. Requires all five toolchains on PATH.
set -euo pipefail

make docs OPAM="opam exec --"
git diff --exit-code
