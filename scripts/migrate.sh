#!/bin/sh

migrate -path db/migrations \
  -database postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=${DB_SSL_MODE} \
  up
