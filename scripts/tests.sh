#!/bin/bash

export GO111MODULE=on 
go fmt $(go list ./... | grep -v /vendor/)
go test -race -coverpkg=./... -v -coverprofile .testCoverage.out ./...
gocov convert .testCoverage.out | gocov report
