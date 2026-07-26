#!/usr/bin/env bash
# Install a pinned ktlint from the official pinterest/ktlint GitHub release.
# Not a piped installer: a versioned binary from a named source, placed on PATH.
set -euo pipefail

ver="${1:-1.8.0}"
url="https://github.com/pinterest/ktlint/releases/download/${ver}/ktlint"
curl -sSL -o ktlint "$url"
chmod +x ktlint
sudo mv ktlint /usr/local/bin/ktlint
ktlint --version
