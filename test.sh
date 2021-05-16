#!/bin/sh
export CGO_ENABLED=0
gofmt -d ./
go test ./... -v -cover | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
