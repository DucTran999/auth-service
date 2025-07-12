#!/bin/bash
set -euo pipefail

mkdir -p test/coverage

# Run unit tests
go test -cover ./internal/usecase/... ./internal/repository/... -coverprofile=test/coverage/unit.out
