#!/usr/bin/env bash
# shellcheck disable=SC1091

set -euo pipefail

PROJECT_DIR="$(dirname "$(dirname "$(realpath "$0")")")"

CYBER_CACHE="${PROJECT_DIR}/.cyber.sh"
CYBER_URL="https://raw.githubusercontent.com/Noksa/install-scripts/main/cyberpunk.sh"

if [[ ! -f "${CYBER_CACHE}" ]]; then
    curl -s "${CYBER_URL}" > "${CYBER_CACHE}"
fi

source "${CYBER_CACHE}"

trap cyber_trap SIGINT SIGTERM
