#!/usr/bin/env sh
set -e
echo "Building..."

BIN_OUT=./bin/logs-converter

go build -o ${BIN_OUT} ./cmd/logs-converter-cli
echo "Binary compiled at ${BIN_OUT}"