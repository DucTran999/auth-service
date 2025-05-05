#!/bin/bash

ENV_FILE_PATH=.env
PKG_INFRA=./test/dependencies
docker-compose -f $PKG_INFRA/docker-compose-test.yml --env-file $ENV_FILE_PATH up -d