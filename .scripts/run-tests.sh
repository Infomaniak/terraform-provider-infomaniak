#!/bin/bash
set -e

REPORT_DIR=$1
export TF_TESTS_MOCKED="true"

echo "=== Running acceptance mocked tests with coverage ==="
GOCOVERDIR="$REPORT_DIR/covcounters" go test ./...
go tool covdata merge -i="$REPORT_DIR/covcounters" -o="$REPORT_DIR/covcounters/merged"
go tool covdata textfmt -i="$REPORT_DIR/covcounters" -o="$REPORT_DIR/acceptance/coverage.out"
echo "Exported covcounters at $REPORT_DIR/covcounters"

echo "=== Running unit mocked tests with coverage ==="
go test -v -coverprofile="$REPORT_DIR/unit/coverage.out" -json ./... > "$REPORT_DIR/unit/results.json"
echo "Test results saved to $REPORT_DIR/unit/results.json"
