#!/bin/bash
set -e

REPORT_DIR=$1
OS=$2

COV_FILE_UNIT="$REPORT_DIR/unit/coverage.out"

if [[ -f "$COV_FILE_UNIT" ]]; then
    awk '/^total:/ {gsub(/%/, "", $3); print $3}' "$REPORT_DIR/unit/coverage_summary.txt" > "$REPORT_DIR/unit/global_coverage.txt"

    UNIT_COVERAGE=$(cat "$REPORT_DIR/unit/global_coverage.txt")

    echo "Global coverage unit: ${UNIT_COVERAGE}%"
else
    echo "0" > "$REPORT_DIR/unit/global_coverage.txt"
    echo "No coverage data found"
fi


COV_ACCEPTANCE_MOCKED="$REPORT_DIR/acceptance/coverage.out"

if [[ -f "$COV_ACCEPTANCE_MOCKED" ]]; then
    if [[ "$OS" == "Darwin" ]]; then
        sed -i '' 's|.*/terraform-provider-infomaniak/|terraform-provider-infomaniak/|g' "$REPORT_DIR/acceptance/coverage_summary.txt"
    else
        sed -i 's|.*/terraform-provider-infomaniak/|terraform-provider-infomaniak/|g' "$REPORT_DIR/acceptance/coverage_summary.txt"
    fi

    awk '/^total:/ {gsub(/%/, "", $3); print $3}' "$REPORT_DIR/acceptance/coverage_summary.txt" > "$REPORT_DIR/acceptance/global_coverage.txt"
    ACC_COVERAGE=$(cat "$REPORT_DIR/acceptance/global_coverage.txt")
    echo "Global coverage acceptance: ${ACC_COVERAGE}%"
else
    echo "0" > "$REPORT_DIR/acceptance/global_coverage.txt"
    echo "No coverage data found"
fi
