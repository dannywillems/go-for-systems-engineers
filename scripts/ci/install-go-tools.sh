#!/usr/bin/env bash
# Install the Go static-analysis tools the lint targets require.
set -euo pipefail

go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
