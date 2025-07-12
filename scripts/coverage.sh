#!/bin/bash
set -euo pipefail

mkdir -p test/coverage

gocovmerge test/coverage/unit.out test/coverage/integration.out >test/coverage/coverage.out

# Optional: generate HTML report
go tool cover -html=test/coverage/coverage.out -o test/coverage/coverage.html

# Calculate total coverage
total_coverage=$(go tool cover -func=test/coverage/coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
coverage_threshold=80.0
comparison=$(echo "$total_coverage >= $coverage_threshold" | bc -l)

if [ "$comparison" -eq 0 ]; then
  echo -e "\033[31m❌ Total coverage: ${total_coverage}% is below threshold (${coverage_threshold}%)\033[0m"
  exit 1
else
  echo -e "\033[32m✅ Total coverage: ${total_coverage}% meets threshold (${coverage_threshold}%)\033[0m"
fi
