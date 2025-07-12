#!/bin/bash
set -euo pipefail

if [ -f .test.env ]; then
    export $(grep -v '^#' .test.env | xargs)
fi

mkdir -p ./test/coverage

# 🧪 Run integration tests
go test -coverpkg=./internal/repository/... -coverprofile=test/coverage/integration.out ./test/integration/...
