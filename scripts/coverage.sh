#!/bin/bash
set -euo pipefail

mkdir -p test/coverage

# Tooling checks
command -v gocovmerge >/dev/null 2>&1 || {
  echo "gocovmerge not found. Install with: go install github.com/wadey/gocovmerge@latest" >&2
  exit 1
}
command -v bc >/dev/null 2>&1 || {
  echo "bc not found. Install via your package manager (e.g., apt-get install bc)" >&2
  exit 1
}

# Input file checks
for f in test/coverage/unit.out test/coverage/integration.out; do
  [[ -f $f ]] || {
    echo "$f not found. Run the corresponding tests to generate this file." >&2
    exit 1
  }
done

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
