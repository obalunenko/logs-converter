#!/usr/bin/env bash
echo "metalinter..."
gometalinter --vendor ./...
echo "go vet..."
go vet ./...
echo "golint..."
golint $(go list ./... | grep -v /vendor/)
echo "gogroup..."
gogroup -order std,other,prefix=github.com/oleg.balunenko/  $(find . -type f -name "*.go" | grep -v "vendor/")