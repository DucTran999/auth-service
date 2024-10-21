#!/bin/bash
go test -v -race $(go list ./... | grep -E "internal/(repo|service|handler)")