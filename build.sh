#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

mkdir -p out

echo "Building commands"
export GOOS=linux
export GOARCH=amd64

go clean
go build -o out/check   cmd/check/main.go
go build -o out/in      cmd/in/main.go

echo "Building and tagging docker image"
docker build . --tag gcr.io/cf-elafros-dog/knative-service-resource

echo "Done."
