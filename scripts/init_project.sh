#!/bin/bash

ENV_FILE_PATH=".env"
if [ ! -e "$ENV_FILE_PATH" ]; then
    cp .env.example .env
fi

DOCKER_COMPOSE_FILE="infra/docker-compose.yml"
if [ ! -e "$DOCKER_COMPOSE_FILE" ]; then
    cp infra/docker-compose-example.yml infra/docker-compose.yml
fi

go mod tidy