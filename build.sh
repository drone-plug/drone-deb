#!/bin/sh

set -e

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build -t plugins/deb .
