#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

if ! command -v goimports &>/dev/null; then
    echo "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
fi

go mod tidy
go fmt ./...
goimports -w .
go vet ./...
go run golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@latest -fix ./...
golangci-lint run

echo "✓ All checks passed"
