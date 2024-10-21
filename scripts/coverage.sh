#!/bin/bash

mkdir -p coverage
go test -cover $(go list ./... | grep -E "internal/(repo|service|handler)")  -coverprofile=coverage/coverage.out

# Enforce > 60% coverage
total_coverage=$(go tool cover -func=coverage/coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
echo "Total coverage: $total_coverage%"
coverage_threshold=60.0
if (( $(echo "$total_coverage < $coverage_threshold" | bc -l) )); then
  echo "Code coverage $total_coverage% is below the threshold of $coverage_threshold%."
  exit 1
else
  echo "Code coverage $total_coverage% meets the threshold of $coverage_threshold%."
fi