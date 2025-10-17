#!/bin/bash
set -e

REPORT_DIR=$1
OS=$2

COV_UNIT="$REPORT_DIR/unit/coverage.out"
COV_ACCEPTANCE="$REPORT_DIR/acceptance/coverage.out"

if [[ -f "$COV_UNIT" ]]; then
    go tool github.com/boumenot/gocover-cobertura < "$COV_UNIT" > "$REPORT_DIR/unit/cobertura.xml"

    echo "Cobertura report generated: $REPORT_DIR/unit/cobertura.xml"
else
    echo "No coverage data found for unit tests"
fi

if [[ -f "$COV_ACCEPTANCE" ]]; then
    if [[ "$OS" == "Darwin" ]]; then
        sed -i '' 's|.*/terraform-provider-infomaniak/|terraform-provider-infomaniak/|g' "$COV_ACCEPTANCE"
    else
        sed -i 's|.*/terraform-provider-infomaniak/|terraform-provider-infomaniak/|g' "$COV_ACCEPTANCE"
    fi

    go tool github.com/boumenot/gocover-cobertura < "$COV_ACCEPTANCE" > "$REPORT_DIR/acceptance/cobertura.xml"

    echo "Cobertura report generated: $REPORT_DIR/acceptance/cobertura.xml"
else
    echo "No coverage data found for mocked acceptance tests"
fi