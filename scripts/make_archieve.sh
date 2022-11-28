#!/usr/bin/env bash

mkdir -p ../generated
cd ../generated
version="$(cat ../plugin.yaml | grep "version" | cut -d '"' -f 2)"
ARCH="amd64 arm64"
OS="linux darwin"
for A in $ARCH; do
  for O in $OS; do
    GOARCH=$A GOOS=$O go build -o "resolve-deps" ../cmd/resolve_deps.go
    tar -czvf "resolve-deps_${version}_${O}_${A}.tar.gz" "resolve-deps"
    rm -rf "resolve-deps"
  done
done