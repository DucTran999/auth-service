#!/bin/bash
set -euo pipefail

mkdir -p test/coverage

# Run unit tests
go test -covermode=set -coverprofile=test/coverage/unit.out \
    ./internal/usecase/... \
    ./internal/repository/... \
    ./internal/handler/rest/...
