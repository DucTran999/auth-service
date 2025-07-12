#!/usr/bin/env bash

if [ -f .test.env ]; then
    set -a
    source .test.env
    set +a
fi

docker compose -f environment/docker-compose.yml up -d db

docker compose -f environment/docker-compose.yml up -d redis

docker compose -f environment/docker-compose.yml run --rm migrate
