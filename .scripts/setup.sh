#!/bin/bash
set -e

LOG_DIRS=$1
REPORT_DIR=$2

for dir in $LOG_DIRS; do
    mkdir -p "$dir"
    touch "$dir/terraform.log"
done

mkdir -p "$REPORT_DIR"/{unit,acceptance,covcounters/merged}