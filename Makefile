#!/usr/bin/make -f

test: fmt
	GORACE="atexit_sleep_ms=50" go test -timeout=1s -race -count 1 -covermode=atomic ./...
	GORACE="atexit_sleep_ms=50" go test -count=1 github.com/smarty/sqldb/integration

fmt:
	go mod tidy && go fmt ./...

compile:
	go build ./...

build: test compile

.PHONY: test fmt compile build
