#!/usr/bin/env bash

function resolve_deps() {
  rm -rf "${1}/charts" "${1}/tmpcharts" "${1}/Chart.lock"
  for dep in $(helm dep list "${1}" | grep "file://" | cut -f 3 | sed s#file:/#.#); do
    resolve_deps "${dep}"
  done

  echo "Running 'helm dep up' in '${1}'"
  helm dep up "${1}"

  if [[ -d "${1}"/charts ]]; then
    for archive in $(find "${1}/charts" -maxdepth 1 -name "*.tgz"); do
      tar xzf "${archive}" -C "${1}/charts"
      rm -rf "${archive}"
    done
  fi
}

if ! command -v tar &>/dev/null; then
  echo "tar program wasn't found. Install it first." >&2
  return 1
fi

echo "Resolving dependencies..."
resolve_deps "${1}"
