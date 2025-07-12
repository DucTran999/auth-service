#!/usr/bin/env bash

docker compose -f environment/docker-compose.yml --env-file .env up -d db

docker compose -f environment/docker-compose.yml --env-file .env up -d redis

docker compose -f environment/docker-compose.yml --env-file .env run --rm migrate
