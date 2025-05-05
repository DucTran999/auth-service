#!/bin/bash

ENV_FILE_PATH=.env.test
PKG_INFRA=./infra
docker-compose -f $PKG_INFRA/docker-compose-test.yml --env-file $ENV_FILE_PATH up -d