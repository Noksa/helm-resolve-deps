#!/usr/bin/env bash

mkdir -p ../generated
cd ../generated
version="$(cat ../plugin.yaml | grep "version" | cut -d '"' -f 2)"
ARCH="amd64 arm64"
OS="linux darwin windows"
for A in $ARCH; do
  for O in $OS; do
    output="resolve-deps"
    if [[ "$O" == "windows" ]]; then
      output="resolve-deps.exe"
    fi
    GOARCH=$A GOOS=$O go build -ldflags "-X main.version=${version}" -o "${output}" ../cmd/resolve_deps.go
    tar -czvf "resolve-deps_${version}_${O}_${A}.tar.gz" "${output}"
    #rm -rf "${output}"
  done
done