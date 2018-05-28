#!/usr/bin/env bash

set -e

echo "Building Waffle Brain bin..."
cd ./brain
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
chmod +x main
echo "Building Waffle Proxy bin..."
cd ../proxy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
chmod +x main
echo "Building Waffle CLI..."
cd ../cli
go build -o waffle-cli main.go
chmod +x waffle-cli

echo "Building Docker image for Waffle Proxy..."
cd ../proxy
docker build -t waffle.io/waffle-proxy:latest .
rm main
cd ../brain
echo "Building Docker image for Waffle Brain..."
docker build -t waffle.io/waffle-brain:latest .
rm main
cd ../proxy-init
echo "Building Docker image for Waffle Proxy sidecar-init..."
docker build -t waffle.io/waffle-proxy-init:latest .

echo "Build images OK! Go Go Go!"