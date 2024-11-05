#!/bin/bash

ENV_FILE_PATH=".env"
if [ ! -e "$ENV_FILE_PATH" ]; then
    cp .env.example .env
fi

DOCKER_COMPOSE_FILE="infra/docker-compose.yml"
if [ ! -e "$DOCKER_COMPOSE_FILE" ]; then
    cp infra/docker-compose-example.yml infra/docker-compose.yml
fi

# install go migrate tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

go mod tidy