#!/usr/bin/env bash

set -e

protoc -I api --go_out=plugins=grpc:api/gen api/*.proto