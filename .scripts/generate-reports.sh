#!/bin/bash
set -e

REPORT_DIR=$1
OS=$2

COV_FILE_UNIT="$REPORT_DIR/unit/coverage.out"

if [[ -f "$COV_FILE_UNIT" ]]; then
    echo "Generating coverage unit tests reports..."
    
    go tool cover -func="$COV_FILE_UNIT" > "$REPORT_DIR/unit/coverage_summary.txt"
    go tool cover -html="$COV_FILE_UNIT" -o "$REPORT_DIR/unit/coverage.html"

    cp $REPORT_DIR/unit/coverage.html $REPORT_DIR/deploy/coverage_unit.html

    echo "Reports:"
    echo " - Summary: $REPORT_DIR/unit/coverage_summary.txt"
    echo " - HTML: $REPORT_DIR/unit/coverage.html"
else
    echo "No coverage data found for unit tests."
fi

COV_ACCEPTANCE_MOCKED="$REPORT_DIR/acceptance/coverage.out"

if [[ -f "$COV_FILE_UNIT" ]]; then
    echo "Generating coverage mocked acceptance test reports..."
    
    go tool cover -func="$COV_ACCEPTANCE_MOCKED" > "$REPORT_DIR/acceptance/coverage_summary.txt"
    go tool cover -html="$COV_ACCEPTANCE_MOCKED" -o "$REPORT_DIR/acceptance/coverage.html"

    cp $REPORT_DIR/acceptance/coverage.html $REPORT_DIR/deploy/coverage_mocked_acceptance.html

    echo "Reports:"
    echo " - Summary: $REPORT_DIR/acceptance/coverage_summary.txt"
    echo " - HTML: $REPORT_DIR/acceptance/coverage.html"
else
    echo "No coverage data found for acceptance mocked tests."
fi
