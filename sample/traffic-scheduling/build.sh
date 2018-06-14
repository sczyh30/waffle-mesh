#!/usr/bin/env bash

set -e

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o foo/main foo/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o baz_v1/main baz_v1/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o baz_v2/main baz_v2/main.go

cd foo
docker build -t waffle.io/sample-ts-foo:latest .
rm main
cd ../baz_v1
docker build -t waffle.io/sample-ts-baz-v1:latest .
rm main
cd ../baz_v2
docker build -t waffle.io/sample-ts-baz-v2:latest .
rm main

echo "Build finished"