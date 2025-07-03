#!/bin/bash

rm -r gen/grpc

buf generate

go mod tidy
