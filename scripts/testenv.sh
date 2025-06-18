#!/bin/sh
set -euo pipefail

docker compose -f environment/docker-compose.yml --env-file .env up -d db

docker compose -f environment/docker-compose.yml run --rm migrate
