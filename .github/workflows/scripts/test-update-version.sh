#!/usr/bin/env bats

setup() {
  source "$BATS_TEST_DIRNAME/update-version.sh"
  
  # Create a temporary directory and dummy Chart.yaml for testing
  TEST_DIR=$(mktemp -d)
  export CHART_PATH="$TEST_DIR/Chart.yaml"
  echo "appVersion: 1.0.0" > "$CHART_PATH"
}

teardown() {
  rm -rf "$TEST_DIR"
}

@test "Use release tag when provided" {
  run determine_version "v1.2.3" "abcdef123456"
  [ "$status" -eq 0 ]
  [ "$output" = "v1.2.3" ]
}

@test "Fall back to short SHA when release tag is missing" {
  run determine_version "" "abcdef123456"
  [ "$status" -eq 0 ]
  [ "$output" = "abcdef1" ]
}

@test "Too short SHA is not extended" {
  run determine_version "" "abc"
  [ "$status" -eq 0 ]
  [ "$output" = "abc" ]
}

@test "Fail when both inputs are missing" {
  run determine_version "" ""
  [ "$status" -eq 1 ]
  [[ "$output" =~ "ERROR: No Release Tag" ]]
}

@test "Full execution updates the file successfully" {
  VERSION_TO_SET="v2.0.0"

  # Ensure file is not yet updated
  [[ ! "$(cat "$CHART_PATH")" =~ "$VERSION_TO_SET" ]]

  run main "$VERSION_TO_SET" ""
  [ "$status" -eq 0 ]
  [[ "$(cat "$CHART_PATH")" =~ "$VERSION_TO_SET" ]]
}
