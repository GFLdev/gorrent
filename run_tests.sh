#!/usr/bin/env bash

[ -f ./tests.log ] && rm ./tests.log
go clean -cache
go test ./... > ./tests.log