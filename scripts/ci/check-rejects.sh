#!/usr/bin/env bash
# Verify every reject-* demo still FAILS TO COMPILE for the documented reason.
#
# Each reject-* directory is a compile-failure showcase (Modules 01,02,03,04,06,
# 08,09,14) plus the fuzz-catches-a-bug demo (Module 11). Its capture `output`
# runs the build and greps for the expected error class; a NON-EMPTY result means
# the demo still fails as the README claims.
#
# These blocks are portable:false, so `make docs-check` skips them, and reject-*
# dirs are excluded from `make build-*` -- so nothing else in CI recompiles them.
# Without this gate a reject that silently started compiling (a language change,
# or a "fix" to the deliberate error) would leave CI green while the README's
# "this does not compile" became false.
set -uo pipefail

# All five toolchains must be reachable. The OCaml reject commands run a bare
# `dune`, so invoke this script UNDER the opam environment -- `opam exec -- bash
# scripts/ci/check-rejects.sh` in CI, or `make check-rejects` locally (which
# wraps it with $(OPAM)). reject-errcheck needs golangci-lint on PATH.
if [ -d /opt/homebrew/opt/openjdk/bin ]; then
	export PATH="/opt/homebrew/opt/openjdk/bin:$PATH"
fi
export GOTOOLCHAIN=auto

cd "$(dirname "$0")/../.." || exit 1

# Emit "<reject-dir>\t<command>" for every capture output whose dir is a reject-*.
list=$(mktemp)
trap 'rm -f "$list"' EXIT
python3 - >"$list" <<'PY'
import json, glob, os


def strip_jsonc(s):
    # drop // line comments like the capture engine's JSONC reader, but never
    # inside a string (so an https:// value is preserved).
    out, i, n, instr, esc = [], 0, len(s), False, False
    while i < n:
        c = s[i]
        if instr:
            out.append(c)
            if esc:
                esc = False
            elif c == '\\':
                esc = True
            elif c == '"':
                instr = False
            i += 1
        elif c == '"':
            instr = True
            out.append(c)
            i += 1
        elif c == '/' and i + 1 < n and s[i + 1] == '/':
            while i < n and s[i] != '\n':
                i += 1
        else:
            out.append(c)
            i += 1
    return ''.join(out)


for cj in sorted(glob.glob('modules/*/capture.json')):
    mod = os.path.dirname(cj)
    m = json.loads(strip_jsonc(open(cj).read()))
    for name, o in (m.get('outputs') or {}).items():
        d = o.get('dir', '')
        if not d.startswith('reject'):
            continue
        cmd = o['cmd']
        inner = cmd[2] if cmd[:2] == ['bash', '-c'] else ' '.join(cmd)
        print(os.path.join(mod, d) + '\t' + inner)
PY

fail=0
checked=0
while IFS=$'\t' read -r dir cmd; do
	[ -z "$dir" ] && continue
	checked=$((checked + 1))
	out=$(cd "$dir" && bash -c "$cmd" 2>/dev/null)
	if [ -z "$out" ]; then
		echo "FAIL  $dir"
		echo "      expected the demo to still fail, but the error was not found."
		echo "      cmd: $cmd"
		fail=1
	else
		printf 'ok    %-42s %s\n' "$dir" "$out"
	fi
done <"$list"

echo "----"
if [ "$fail" -ne 0 ]; then
	echo "check-rejects: one or more reject demos no longer fail as documented."
	exit 1
fi
echo "check-rejects: all $checked reject demos still fail for the documented reason."
