#!/bin/sh
export CGO_ENABLED=0
test -z "$(gofmt -l -d ./)"
mkdir coverage
go test ./... -v -coverprofile=coverage/coverage.txt -covermode=atomic
