#!/usr/bin/env bash

set -euo pipefail

CHART_PATH="${CHART_PATH:-charts/cloudmetrics/Chart.yaml}"

determine_version() {
  local if_release="$1"
  local github_sha="$2"

  if [ -n "$if_release" ]; then
    echo "$if_release"
  elif [ -n "$github_sha" ]; then
    echo "${github_sha:0:7}"
  else
    echo "ERROR: No Release Tag AND no GitHub SHA provided." >&2
    return 1
  fi
}

main() {
  local if_release="${1:-}"
  local github_sha="${2:-}"

  # Get version or catch error
  if ! VERSION=$(determine_version "$if_release" "$github_sha"); then
    echo "$VERSION" >&2
    exit 1
  fi

  if [ ! -f "$CHART_PATH" ]; then
    echo "ERROR: Chart file not found at $CHART_PATH" >&2
    exit 1
  fi

  # Update the Helm chart in-place using yq
  # Note: This command requires the Go version of 'yq' (mikefarah/yq), 
  # NOT the Python version. (Install via 'snap install yq')
  yq -i ".appVersion = \"${VERSION}\"" "$CHART_PATH"
  echo "Successfully updated appVersion to '${VERSION}' in $CHART_PATH"
}

# Only run main if the script is being executed directly, not sourced
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  main "$@"
fi
