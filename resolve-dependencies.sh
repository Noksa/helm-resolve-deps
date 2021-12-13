#!/usr/bin/env bash

function usage() {
    cat <<EOM
    Usage:
    helm resolve-deps [CHART_DIR] [FLAGS]
    If CHART_DIR is empty, the plugin will try to resolve dependencies in the current directory

    FLAGS:
      -u[--unpack-dependencies]       - untar/unpack all (including external) dependent charts. They will be present as directories instead of .tgz archieves inside chartrs/ directory
      -c[--clean]                     - remove charts, tmpcharts directories and Chart.lock file in each chart before running the dependency update command
      --skip-refresh-in name1,name2   - skip fetching updates from helm repositories before running 'helm dep up' in specific charts (pass their names in the argument).
                                        Use ',' as delimiter if you want to specify more than one chart.

      All flags from 'helm dependency update' command can be passed as flags to the plugin's command

    Examples:
      helm resolve-deps /my-chart --skip-refresh --unpack-dependencies
      helm resolve-deps /my-chart --skip-refresh
      helm resolve-deps --skip-refresh-in my-chart1,my-second-chart
EOM
    exit 0
}

trap 'usage' err

function resolve_deps() {
  local CHART_DIR="${1}"
  local CHART_NAME=""
  CHART_NAME=$(helm show chart "${CHART_DIR}" | grep "^name" | awk '{ print $2 }')
  shift
  if [[ ${CLEAN} == true ]]; then
    rm -rf "${CHART_DIR}/charts" "${CHART_DIR}/tmpcharts" "${CHART_DIR}/Chart.lock"
  fi
  for dep in $(helm dep list "${CHART_DIR}" | grep "file://" | cut -f 3 | sed s#file:/#.#); do
    resolve_deps "${CHART_DIR}/${dep}" $@
  done
  echo "Running 'helm dep up' in '${CHART_DIR}'"
  if [[ ${SKIP_REFRESH_ALL_CHARTS} == true ]]; then
    set -- "$@" "--skip-refresh"
  else
    for SKIP_REFRESH_IN_CHART in "${SKIP_REFRESH_IN_CHARTS[@]}"; do
      if [[ "${CHART_NAME}" == "${SKIP_REFRESH_IN_CHART}" ]]; then
        set -- "$@" "--skip-refresh"
        break
      fi
    done
  fi
  helm dep up "${CHART_DIR}" $@

  if [[ -d "${CHART_DIR}"/charts && ${UNTAR_CHARTS} == true ]]; then
    for archive in $(find "${CHART_DIR}/charts" -maxdepth 1 -name "*.tgz"); do
      tar xzf "${archive}" -C "${CHART_DIR}/charts"
      rm -rf "${archive}"
    done
  fi
}

MAINCHART_DIR="${1}"
if [[ $MAINCHART_DIR == -* ]]; then
  MAINCHART_DIR=
fi
if [[ -z "${MAINCHART_DIR:-}" ]]; then
  MAINCHART_DIR="$(pwd)"
else
  shift
fi
UNTAR_CHARTS=false
CLEAN=false
SKIP_REFRESH_ALL_CHARTS=false
SKIP_REFRESH_IN_CHARTS=()
ARGUMENTS=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    -u | --unpack-dependencies)
       UNTAR_CHARTS=true
       shift
    ;;
    -h | --help)
       usage
       shift
       exit 0
    ;;
    -c | --clean)
      CLEAN=true
      shift
    ;;
    --skip-refresh)
      SKIP_REFRESH_ALL_CHARTS=true
      shift
    ;;
    --skip-refresh-in)
      for c in $(echo "${2}" | tr "," "\n"); do
        SKIP_REFRESH_IN_CHARTS+=( "${c}" )
      done
      shift 2
    ;;
    *)
      ARGUMENTS="${ARGUMENTS} ${1}"
      shift
    ;;
  esac
done
if [[ $UNTAR_CHARTS == true ]]; then
  if ! command -v tar &>/dev/null; then
    echo "tar program wasn't found. Install it first." >&2
    return 1
  fi
fi

echo "Resolving dependencies..."
resolve_deps "${MAINCHART_DIR}" ${ARGUMENTS}
