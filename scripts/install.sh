#!/usr/bin/env sh

set -e

cd $HELM_PLUGIN_DIR
version="$(cat plugin.yaml | grep "version" | cut -d '"' -f 2)"
echo "Installing resolve-deps v${version} ..."

# Find correct archive name
unameOut="$(uname -s)"

case "${unameOut}" in
    Linux*)     os=linux;;
    Darwin*)    os=darwin;;
#    CYGWIN*)    os=cygwin;;
#    MINGW*)     os=windows;;
    *)          os="UNKNOWN:${unameOut}"
esac

arch=$(uname -m)

if [ "$arch" = "aarch64" ]; then
  arch="arm64"
fi

url="https://github.com/Noksa/helm-resolve-deps/releases/download/v${version}/resolve-deps_${version}_${os}_${arch}.tar.gz"

if [ "$url" = "" ]
then
    echo "Unsupported OS / architecture: ${os}_${arch}"
    exit 1
fi

filename="resolve-deps_${version}.tar.gz"


if [ -z "$(command -v tar)" ]; then
  echo "tar is required, install it first"
  exit 1
fi

# Download archive
if [ -n "$(command -v curl)" ]
then
    curl -sSL "$url" -o "$filename"
elif [ -n "$(command -v wget)" ]
then
    wget -q "$url" -o "$filename"
else
    echo "Need curl or wget"
    exit 1
fi

trap 'rm -rf $filename' EXIT

# Install bin
rm -rf bin && mkdir bin && tar xzvf "$filename" -C bin > /dev/null && rm -f "$filename"

if [ "$?" != "0" ]; then
  echo "an error has occured"
  exit 1
fi

echo "resolve-deps ${version} has been installed"
echo
echo "See https://github.com/Noksa/helm-resolve-deps for usage"
