#!/usr/bin/env bash

if [ -f .test.env ]; then
    export $(grep -v '^#' .test.env | xargs)
fi

docker compose -f environment/docker-compose.yml up -d db

docker compose -f environment/docker-compose.yml up -d redis

docker compose -f environment/docker-compose.yml run --rm migrate
