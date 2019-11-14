#!/bin/bash

echo "==Generate Build Script=="
echo "nextRelease.version: $1"
echo "options.branch: $2"
echo "commits.length: $3"
echo "build-date: $4"

go install -ldflags "-X main.Version=$1 -X main.BuildDate=$4" github.com/prijip/gofind/cmd/gofind

gofind --version

docker build -t gofind:latest -f ./semantic-release/Dockerfile $GOPATH/bin

docker run -ti gofind:latest
