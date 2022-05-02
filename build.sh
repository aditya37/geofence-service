#!/usr/bin/env bash

## export go module
export GO111MODULE=on

## export gosumb
export GOSUMDB=off

/usr/local/go/bin/go clean && CGO_ENABLED=0 /usr/local/go/bin/go build

