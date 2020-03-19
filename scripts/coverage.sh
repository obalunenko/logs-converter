#!/usr/bin/env sh

go test -race -coverpkg=./... -v -coverprofile coverage.out ./...
gocov convert coverage.out | gocov report
go tool cover -html=coverage.out
