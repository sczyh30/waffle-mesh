#!/usr/bin/env bash

set -e

rm -rf api/gen
mkdir api/gen

protoc -I api --go_out=plugins=grpc:api/gen api/*.proto